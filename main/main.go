package main

import (
	"flag"
	"k8s.io/klog/v2"
)

func main() {
	var dryRun = flag.Bool("dry-run", false, "dry-run mode - doesn't delete anything")
	klog.InitFlags(flag.CommandLine)
	flag.Parse()

	kubeApi := CreateKubeApi()
	err := kubeApi.CleanPods(*dryRun)
	if err != nil {
		klog.Fatal(err)
	}
}
