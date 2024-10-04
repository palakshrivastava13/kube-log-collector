package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetPodLog collects logs from a particular pod and saves them in a file
func GetPodLog(kubeconfig, namespace, podName, outputDir string) error {
	// Load the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	// Get the logs of the specified pod
	logOptions := &v1.PodLogOptions{}
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)

	logs, err := req.Stream(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to get logs for pod %s: %v", podName, err)
	}
	defer logs.Close()

	// Create a log file for the pod
	logFilePath := filepath.Join(outputDir, podName)
	file, err := CreateLogFile(logFilePath, podName)
	if err != nil {
		return fmt.Errorf("failed to create log file for pod %s: %v", podName, err)
	}
	defer file.Close()

	// Save the logs to the file
	if _, err := io.Copy(file, logs); err != nil {
		return fmt.Errorf("failed to write logs to file for pod %s: %v", podName, err)
	}

	fmt.Printf("Logs for pod %s saved to %s\n", podName, file.Name())
	return nil
}

// GetPodLogs collects logs from all pods in a namespace and saves them in a directory
func GetPodLogs(kubeconfig, namespace, outputDir string) error {
	// Load the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %v", err)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	// List all pods in the specified namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

	// Create output directory for namespace logs
	namespaceDir := filepath.Join(outputDir, namespace)
	if err := os.MkdirAll(namespaceDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for namespace logs: %v", err)
	}

	// Iterate over each pod and get the logs using the GetPodLog method
	for _, pod := range pods.Items {
		if err := GetPodLog(kubeconfig, namespace, pod.Name, namespaceDir); err != nil {
			fmt.Printf("Error collecting logs for pod %s: %v\n", pod.Name, err)
		}
	}

	return nil
}
