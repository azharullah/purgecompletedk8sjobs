# PurgeK8sJobs

[![Go Report Card](https://goreportcard.com/badge/github.com/azharullah/purgek8sjobs)](https://goreportcard.com/report/github.com/azharullah/purgek8sjobs)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/azharullah/purgek8sjobs?tab=doc)

PurgeK8sJobs is package to purge the completed Kubernetes Jobs.

### Why PurgeK8sJobs?

1) As of the time of development of this project, Kubernetes already provides a feature to achieve this, [documented here](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#ttl-mechanism-for-finished-jobs). However, it is still an alpha feature and on cloud providers like GCP, using alpha features is not recommended and the Alpha GKE clusters expire after a month.
2) Moreover, using this package would provide additional (optional) features that help in user debugging and tracking. The features include:
   1) A user configurable parameter - `hrs`, which implies that the K8s jobs that have completed in the last `hrs` hours will not be purged.
   2) Optionally, dumping the job's spec to a log file that the user can provide in the options.
   3) Optionally, dumping the job's events to a log file that the user can provide in the options.

**Note**: The events can be captured unless they have not been purged by the K8s API server. This is configurable via the API server flag - `--event-ttl`, which defaults to 1hr. Hence, with the default API server configuration, if the job is not getting purged within one hour of it's completion, it's events will not be available.

#### Usage

```go
options := map[string]string{
    "eventsLogFile": "./events.txt",
    "specLogFile":   "./specs.txt",
}

resp := purgek8sjobs.PurgeJobs("default", 5, options)
if resp.Success {
    log.Print(resp.Msg)
} else {
    log.Fatalf("Failed to delete some / all the compeleted job(s), error: %v", resp.Err.Error())
}
```

There is a Golang based CLI tool [available here](https://github.com/azharullah/purge-k8s-jobs-cli), that provides a CLI interface to this package.

#### Kubernetes authentication

The package uses the [client-go](https://github.com/kubernetes/client-go) library internally to authenticate to the Kubernetes API server. 
And right now supports two ways of authentication that the client-go library supports. The package first tries to authenticate using the method 1. If that fails, it falls back to the method 2 of authentication.

1) [In-Cluster client configuration](https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration): This is useful when this package is being used in a program that runs inside a pod and needs to authenticate to the cluster using the service accounts' token and certificate mounted to the pod. Care should be taken that the service account mounted has enough RBAC permissions to list and delete the BatchV1 Jobs in the namespace being used.
2) [Out-of-Cluster client configuration](https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration): This is useful for manual debugging or when this package is being used in programs that run outside of the cluster. This reads the `KUBECONFIG` environment variable and uses the kubeconfig to authenticate to the cluster.