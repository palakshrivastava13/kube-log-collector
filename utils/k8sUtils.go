package utils

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os/exec"
)

// FetchKubeconfig Fetch the kubeconfig based on cloud type and cluster name
func FetchKubeconfig(cloudType, clusterName, authDetails, kubeconfigPath string) error {
	switch cloudType {
	case "AWS":
		cmd := exec.Command("aws", "eks", "update-kubeconfig", "--name", clusterName, "--kubeconfig", kubeconfigPath, "--profile", authDetails)
		return cmd.Run()
	case "GCP":
		cmd := exec.Command("gcloud", "container", "clusters", "get-credentials", clusterName, "--zone", authDetails, "--kubeconfig", kubeconfigPath)
		return cmd.Run()
	case "Azure":
		cmd := exec.Command("az", "aks", "get-credentials", "--name", clusterName, "--kubeconfig", kubeconfigPath)
		return cmd.Run()
	default:
		return fmt.Errorf("unsupported cloud type: %s", cloudType)
	}
}

// GetKubeClient Get Kubernetes clientset
func GetKubeClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}
	return clientset, nil
}
