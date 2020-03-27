package purgek8sjobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// logJobSpecToFile fetches the provided job's spec and converts the JSON spec to a formatted
// string and then appends it to the provided file. Raises an error if the spec cannot be
// appended to a file
func logJobSpecToFile(j batchv1.Job, f string) (err error) {

	jobSpecBytes, err := json.MarshalIndent(j.Spec, "", "\t")
	if err != nil {
		return errors.New("Failed to write data to config, error: " + err.Error())
	}

	logMsg := fmt.Sprintf("\nSpec for job [%v], which was deleted at %v", j.GetName(), time.Now().String())
	err = appendBytesToFile(jobSpecBytes, f, logMsg)
	if err != nil {
		return err
	}

	return nil
}

// logJobEventsToFile takes the jobs event list and then appends it to the provided file.
// Raises an error if the spec cannot be appended to a file
func logJobEventsToFile(j batchv1.Job, je corev1.EventList, f string) (err error) {

	jobEventBytes, err := json.MarshalIndent(je, "", "\t")
	if err != nil {
		return errors.New("Failed to marshal data, error: " + err.Error())
	}

	logMsg := fmt.Sprintf("\nEvents for job [%v], which was deleted at %v", j.GetName(), time.Now().String())
	err = appendBytesToFile(jobEventBytes, f, logMsg)
	if err != nil {
		return err
	}
	return nil
}

// appendBytesToFile takes a byte slice and an optional message and
// appends them to the provided file - `f`
func appendBytesToFile(b []byte, f string, m string) (err error) {

	sep := strings.Repeat("*", 80)

	fullLogStr := []byte("\n")
	fullLogStr = append(fullLogStr, sep...)
	fullLogStr = append(fullLogStr, m...)
	fullLogStr = append(fullLogStr, "\n"...)
	fullLogStr = append(fullLogStr, b...)
	fullLogStr = append(fullLogStr, "\n"...)
	fullLogStr = append(fullLogStr, sep...)

	// Get the absolute path of the log file
	fileAbsPath, err := homedir.Expand(f)
	if err != nil {
		return errors.New("Failed to expand log file path, error: " + err.Error())
	}

	// Open the log file in append mode, create if it does not exist
	fileHandler, err := os.OpenFile(fileAbsPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("Failed to open the file to write content, error: " + err.Error())
	}
	defer fileHandler.Close()

	// Append the content to the file
	logrus.Debugf("Writing content to file: %v", fileHandler.Name())
	_, err = fileHandler.Write([]byte(fullLogStr))
	if err != nil {
		return errors.New("Failed to write the content to the file, error: " + err.Error())
	}

	return nil
}
