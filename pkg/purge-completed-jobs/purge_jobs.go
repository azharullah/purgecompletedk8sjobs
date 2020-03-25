package purgecompletedk8sjobs

import (
	"log"
	"time"
)

// PurgeResponse : PurgeResponse
type PurgeResponse struct {
	Success bool
	Msg     string
	Err     error
}

// PurgeJobs : PurgeJobs
func PurgeJobs(ns string, hrs int16, options map[string]string) PurgeResponse {

	// Create the k8s client object
	kClient := getK8sAPIClient()

	// Compute the time before `hrs` hours
	currentTime := time.Now()
	reqTime := currentTime.Add(time.Duration(-hrs) * time.Hour)

	log.Printf("Will attempt to delete all Jobs that got completed before %v (ie. before %d hrs)", reqTime, hrs)

	// Get all the jobs that have completed before reqTime
	jobsToDelete, err := getEligibleJobs(kClient, ns, reqTime)
	if err != nil {
		return PurgeResponse{Success: false, Err: err}
	}

	noOfJobs := len(jobsToDelete)
	if noOfJobs > 0 {

		// If there are any such jobs, delete them
		log.Printf("Found %d jobs to delete, will attempt to delete them", noOfJobs)
		msg, err := deleteJobs(kClient, ns, jobsToDelete)
		if err != nil {
			return PurgeResponse{Success: false, Err: err}
		}
		return PurgeResponse{Success: true, Msg: msg}

	}

	log.Printf("Found no eligible jobs to delete, returning...")
	return PurgeResponse{
		Success: true,
		Msg:     "Found no jobs to delete",
		Err:     nil,
	}
}
