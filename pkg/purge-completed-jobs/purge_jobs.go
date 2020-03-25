package purgecompletedk8sjobs

import (
	"errors"
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

	currentTime := time.Now()
	reqTime := currentTime.Add(time.Duration(-hrs) * time.Hour)

	log.Printf("Will attempt to delete all Jobs that got completed before %v (ie. before %d hrs)", reqTime, hrs)

	kClient := getK8sAPIClient()

	jobsToDelete, err := getEligibleJobs(kClient, ns, reqTime)
	if err != nil {
		return PurgeResponse{
			Success: false,
			Err:     err,
		}
	}

	noOfJobs := len(jobsToDelete)
	if noOfJobs > 0 {

		log.Printf("Found %d jobs to delete, will attempt to delete them", noOfJobs)
		// delete jobs

	} else {

		log.Printf("Found no eligible jobs to delete, returning...")
		return PurgeResponse{
			Success: true,
			Msg:     "Found no jobs to delete",
			Err:     nil,
		}

	}

	return PurgeResponse{
		Success: false,
		Err:     errors.New("Failed to process / delete the completed jobs"),
	}
}
