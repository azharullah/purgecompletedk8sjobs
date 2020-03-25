package purgecompletedk8sjobs

import (
	"log"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getEligibleJobs(c *kubernetes.Clientset, n string, t time.Time) (jobList []batchv1.Job, err error) {
	jobList = []batchv1.Job{}

	// Fetch all the k8s jobs in `n` namespace
	listOptions := metav1.ListOptions{
		// LabelSelector: "{.items[?(@.status.completionTime<=\"2020-03-14T08:12:57Z\")]}",
	}
	allJobs, err := c.BatchV1().Jobs(n).List(listOptions)

	for _, job := range allJobs.Items {

		if job.Status.Active == 0 { // Filter out non-active jobs

			kubeTimeObj := metav1.NewTime(t) // Convert the time obj to K8s time obj for comparision

			if job.Status.CompletionTime.Before(&kubeTimeObj) { // Filter out jobs that finished before t
				log.Printf("Got an eligible job with name: %#v", job.GetName())
				jobList = append(jobList, job)
			}
		}
	}

	return jobList, nil

}

// propagationPolicy := metav1.DeletePropagationForeground
// deleteOptions := &metav1.DeleteOptions{
// 	PropagationPolicy: &propagationPolicy,
// }
// err = c.BatchV1().Jobs(n).Delete(job.GetName(), deleteOptions)
// if err != nil {
// 	log.Fatalf("Failed to delete the job, error: %v", err.Error())
// }
