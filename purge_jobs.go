// Package purgecompletedk8sjobs provides a method to delete K8s jobs which have been
// completed, with a few additional and optional features to help with debugging and
// trackability.
package purgecompletedk8sjobs

import (
	"time"

	"github.com/sirupsen/logrus"
)

// PurgeResponse captures the standard for the response of the PurgeJobs function
type PurgeResponse struct {
	Success bool
	Msg     string
	Err     error
}

// PurgeJobs : PurgeJobs
func PurgeJobs(ns string, hrs int16, options map[string]string) PurgeResponse {

	// Create the k8s client object
	logrus.Debug("Creating the authenticated K8s client")
	kClient := getK8sAPIClient()

	// Compute the time before `hrs` hours
	currentTime := time.Now()
	reqTime := currentTime.Add(time.Duration(-hrs) * time.Hour)

	logrus.WithFields(logrus.Fields{
		"before-hrs":  hrs,
		"before-time": reqTime.String(),
	}).Info("Will attempt to delete all Jobs that got completed")

	// Get all the jobs that have completed before reqTime
	logrus.Debug("Getting the eligible jobs to be deleted")
	jobsToDelete, err := getEligibleJobs(kClient, ns, reqTime)
	if err != nil {
		return PurgeResponse{Success: false, Err: err}
	}

	noOfJobs := len(jobsToDelete)
	if noOfJobs > 0 {

		// If there are any such jobs, delete them
		logrus.Infof("Found %d jobs to delete, will attempt to delete them", noOfJobs)
		msg, err := deleteJobs(kClient, ns, jobsToDelete, options)
		if err != nil {
			return PurgeResponse{Success: false, Err: err}
		}
		return PurgeResponse{Success: true, Msg: msg}
	}

	logrus.Debug("Found no eligible jobs to delete, returning...")
	return PurgeResponse{
		Success: true,
		Msg:     "Found no jobs to delete",
		Err:     nil,
	}
}
