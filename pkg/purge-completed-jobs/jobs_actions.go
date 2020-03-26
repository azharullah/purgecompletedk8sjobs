package purgecompletedk8sjobs

import (
	"errors"
	"log"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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
				log.Printf("Got an eligible job with name: %v", job.GetName())
				jobList = append(jobList, job)
			}
		}
	}

	return jobList, nil

}

func deleteJobs(c *kubernetes.Clientset, ns string, jl []batchv1.Job, op map[string]string) (retMsg string, err error) {

	retMsg = ""
	deleteSuccesses, deleteFailures := []string{}, []string{}

	// Set the propagation Policy to Foreground
	propagationPolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &propagationPolicy,
	}

	specLogFile, writeSpec := op["specLogFile"]
	eventsLogFile, writeEvents := op["eventsLogFile"]

	for _, job := range jl {

		jobName := job.GetName()

		// Write the job spec to the log file, if provided
		if writeSpec {
			log.Printf("Attempting to write spec to file for job: %v", jobName)
			logJobSpecToFile(job, specLogFile)
		}

		// Write the job events to the log file, if provided
		if writeEvents {
			jobEvents, err := getJobEvents(c, ns, job)
			if err != nil {
				log.Printf("Failed to get events for the job [%v], error: %v", jobName, err.Error())
			} else {
				log.Printf("Attempting to write events to file for job: %v", jobName)
				err = logJobEventsToFile(job, *jobEvents, eventsLogFile)
				if err != nil {
					log.Printf("Failed to write events for the job [%v], error: %v", jobName, err.Error())
				}
			}
		}

		// os.Exit(1)

		err = c.BatchV1().Jobs(ns).Delete(jobName, deleteOptions)
		if err != nil {
			deleteFailures = append(deleteFailures, "Failed to delete job [%v], error: %v", jobName, err.Error())
		}
		deleteSuccesses = append(deleteSuccesses, jobName)
	}

	if len(deleteSuccesses) > 0 {
		retMsg += "\nSuccessfully deleted the following jobs:\n"
		retMsg += strings.Join(deleteSuccesses, "\n")
	}

	if len(deleteFailures) > 0 {
		retMsg += "\nFailed to deleted the following jobs:\n"
		retMsg += strings.Join(deleteFailures, "\n")
	}

	return retMsg, nil
}

func getJobEvents(c *kubernetes.Clientset, n string, j batchv1.Job) (eventList *corev1.EventList, err error) {
	eventListOptions := metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + j.GetName(),
	}
	eventList, err = c.CoreV1().Events(n).List(eventListOptions)
	if err != nil {
		return &corev1.EventList{}, errors.New("Failed to query events from the K8s API, error: " + err.Error())
	}
	return eventList, nil
}
