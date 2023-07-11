# K8s Cronjob Pod Cleaner
Due to a bug in k8s old `pods` of cronjobs are not cleaned correctly when the job is removed.

This utility removes any `pods` which do have a reference to a non-existing `Job`.

## Usage
Create a cronjob (yes I understand the irony here)
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: 'pod-cleaner'
spec:
  concurrencyPolicy: Forbid
  failedJobsHistoryLimit: 1
  successfulJobsHistoryLimit: 1
  schedule: '*/30 * * * *'
  jobTemplate:
    metadata:
      name: pod-cleaner
      labels:
        app: pod-cleaner
    spec:
      activeDeadlineSeconds: 240
      ttlSecondsAfterFinished: 120
      template:
        metadata:
          labels:
            app: "pod-cleaner"
        spec:
          securityContext:
            runAsUser: 1000
            runAsGroup: 1000
            fsGroup: 1000
          serviceAccountName: "pod-cleaner-sa"
          restartPolicy: Never
          containers:
            - name: cleaner
              image: ghcr.io/davidgiga1993/cronjob-pod-cleaner:latest
              resources:
                requests:
                  cpu: "50m"
                  memory: "50Mi"
                limits:
                  cpu: "150m"
                  memory: "150Mi"
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: pod-cleaner-sa
  namespace: your-ns

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: pod-cleaner
rules:
  - apiGroups: [ "" ]
    resources: [ "pods" ]
    verbs: [ "get", "watch", "list", "delete" ]

  - apiGroups: [ "batch" ]
    resources: [ "jobs" ]
    verbs: [ "get", "watch", "list" ]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: pod-cleaner
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: pod-cleaner
subjects:
  - kind: ServiceAccount
    name: pod-cleaner-sa
    namespace: your-ns
```

## Options
```
-dry-run: Do not delete any pods, just log
```

## Contribute
Since this is hopefully a temporary issue I don't expect any contribution, but feel free to open MRs / issues if you think something is missing.