package purgecompletedk8sjobs

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getK8sAPIClient() *kubernetes.Clientset {
	config, err := getClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func getInClusterConfig() (rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return rest.Config{}, err
	}
	return *config, nil
}

func getClusterConfig() (*rest.Config, error) {
	var kubeconfig string
	kubeconfig = "/Users/azharullah.shariff/.kube/azhar-cluster-1"

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return config, nil
}
