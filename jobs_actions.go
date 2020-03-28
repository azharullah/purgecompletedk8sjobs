package purgek8sjobs

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// getEligibleJobs scans the namespace `ns`, fetches all the jobs and then filters out the jobs that have
// finished their execution before t
func getEligibleJobs(c *kubernetes.Clientset, ns string, t time.Time) (jobList []batchv1.Job, err error) {
	jobList = []batchv1.Job{}

	// Fetch all the k8s jobs in `n` namespace
	listOptions := metav1.ListOptions{
		// LabelSelector: "{.items[?(@.status.completionTime<=\"2020-03-14T08:12:57Z\")]}",
	}
	logrus.Debugf("Fetching the jobs in the %v namespace", ns)
	allJobs, err := c.BatchV1().Jobs(ns).List(listOptions)

	for _, job := range allJobs.Items {

		if job.Status.Active == 0 { // Filter out non-active jobs
			kubeTimeObj := metav1.NewTime(t) // Convert the time obj to K8s time obj for comparision

			if job.Status.CompletionTime.Before(&kubeTimeObj) { // Filter out jobs that finished before t
				logrus.Debugf("Got an eligible job with name: %v", job.GetName())
				jobList = append(jobList, job)
			}
		}
	}

	return jobList, nil
}

// deleteJobs takes in a list of K8s Jobs in the namespace `ns` and deletes them. Optionally, it can dump the job's
// events and spec to log files provided in the options - `op`
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
			logrus.Debugf("Attempting to write spec to file for job: %v", jobName)
			logJobSpecToFile(job, specLogFile)
		}

		// Write the job events to the log file, if provided
		if writeEvents {
			jobEvents, err := getJobEvents(c, ns, job)
			if err != nil {
				logrus.Warnf("Failed to get events for the job [%v], error: %v", jobName, err.Error())
			} else {
				logrus.Debugf("Attempting to write events to file for job: %v", jobName)
				err = logJobEventsToFile(job, *jobEvents, eventsLogFile)
				if err != nil {
					logrus.Warnf("Failed to write events for the job [%v], error: %v", jobName, err.Error())
				}
			}
		}

		// Delete the actual job
		logrus.WithField("job-name", jobName).Info("Deleting the Job with")
		err = c.BatchV1().Jobs(ns).Delete(jobName, deleteOptions)
		if err != nil {
			deleteFailures = append(deleteFailures, "Failed to delete job [%v], error: %v", jobName, err.Error())
		}
		deleteSuccesses = append(deleteSuccesses, jobName)
	}

	// Form the return response summary
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

// getJobEvents takes the K8s job object and quueries the Events API for the job's corresponding events.
// The job's event retaining period depends on the K8s API configuration and is 1 hour by default.
// So, if this is called after an hour with the default API server configuration, there will be no events.
func getJobEvents(c *kubernetes.Clientset, n string, j batchv1.Job) (eventList *corev1.EventList, err error) {
	eventListOptions := metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + j.GetName(),
	}
	logrus.Debugf("Fetching the event for the job: %v", j.GetName())
	eventList, err = c.CoreV1().Events(n).List(eventListOptions)
	if err != nil {
		return &corev1.EventList{}, errors.New("Failed to query events from the K8s API, error: " + err.Error())
	}
	return eventList, nil
}

func init() {

	// Set the log level, if not specified via env
	lvl, ok := os.LookupEnv("LOG_LEVEL")

	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}

	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}

	// set global log level
	logrus.SetLevel(ll)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
}
