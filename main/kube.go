package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/pager"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
)

type KubeApi struct {
	kubeClient *kubernetes.Clientset
	ctx        context.Context
}

func CreateKubeApi() KubeApi {
	kubeConfig := ctrl.GetConfigOrDie()
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Panic(err)
	}
	return KubeApi{
		kubeClient: kubeClient,
		ctx:        context.Background(),
	}
}

func (k *KubeApi) CleanPods(dryRun bool) error {
	core := k.kubeClient.CoreV1()
	err := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		return core.Pods("").List(k.ctx, opts)
	})).EachListItem(k.ctx, metav1.ListOptions{}, func(obj runtime.Object) error {
		pod := obj.(*v1.Pod)

		for _, reference := range pod.OwnerReferences {
			if reference.Kind != "Job" {
				continue
			}
			jobName := reference.Name
			namespace := pod.Namespace

			exists, err := k.JobExists(jobName, namespace)
			if err != nil {
				klog.Errorf("could not get job %v in ns %v", jobName, namespace)
				continue
			}
			if !exists {
				// Delete pod
				if dryRun {
					klog.Infof("would delete pod %v in ns %v", pod.Name, namespace)
					continue
				}

				klog.Infof("delete pod %v in ns %v", pod.Name, namespace)
				err := core.Pods(namespace).Delete(k.ctx, pod.Name, metav1.DeleteOptions{})
				if err != nil {
					klog.Errorf("could not delete pod %v in ns %v : %v", pod.Name, namespace, err)
				}
			}
		}
		return nil
	})
	return err
}

func (k *KubeApi) JobExists(name string, namespace string) (bool, error) {
	job, err := k.kubeClient.BatchV1().Jobs(namespace).Get(k.ctx, name, metav1.GetOptions{})
	if err != nil {
		statErr, ok := err.(*errors.StatusError)
		if ok && statErr.Status().Code == 404 {
			return false, nil
		}
	}
	return job != nil, err
}
