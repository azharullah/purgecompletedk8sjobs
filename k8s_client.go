package purgecompletedk8sjobs

import (
	"os"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getK8sAPIClient authenticates to the kubernetes cluster
// and returns the authenticated client. It first tries the in-cluster
// config way of authentication and then falls back to out-of-cluster
// authentication way.
func getK8sAPIClient() *kubernetes.Clientset {
	config, err := getInClusterConfig()
	if err != nil {
		logrus.Warnf("Failed to use in-cluster configuration, will "+
			"attempt out-of-cluster authentication. Error: %v", err.Error())
		config, err = getClusterConfig()
		if err != nil {
			logrus.Fatalf("Failed to use the out-of-cluster "+
				"authentication too, error: %v", err.Error())
		}
	}

	// Create the client object from in-cluster configuration
	logrus.Debug("Creating the k8s client object")
	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		logrus.Fatalf("Failed to create the K8s client "+
			"authentication too, error: %v", err.Error())
	}
	return clientset
}

// getInClusterConfig uses the service account tokens and cert set
// in the pod volume for authentication
// Recommended for use in k8s cronjobs / jobs
func getInClusterConfig() (rest.Config, error) {
	logrus.Debug("Attempting to authenticate using the In-cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		return rest.Config{}, err
	}
	return *config, nil
}

// getClusterConfig uses the KUBECONFIG set via the env var to authenticate
// Recommended for use in local testing / external scripts
func getClusterConfig() (rest.Config, error) {
	kubeconfig, exists := os.LookupEnv("KUBECONFIG")
	if exists != true {
		logrus.Warnf("KUBECONFIG env var not set")
	}
	logrus.WithField("KUBECONFIG", kubeconfig).Debug("Using kubeconfig from env:")

	// use the current context in kubeconfig
	logrus.Debug("Using the provided KUBECONFIG for the Out-of-cluster config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return rest.Config{}, err
	}

	return *config, nil
}
