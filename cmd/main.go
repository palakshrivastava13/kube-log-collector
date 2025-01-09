package main

import (
	"flag"
	"fmt"
	"kube-log-collector/utils"
	"os"
	"path/filepath"
)

func main() {
	// Define command-line flags for kubeconfig file, namespace, pod name, and output directory
	cloudType := flag.String("cloudType", "aws", "The cloudType for the cluster to collect pod logs from")
	clusterName := flag.String("clusterName", "blue", "The cluster to collect pod logs from")
	authDetails := flag.String("authDetails", "my-auth-profile", "AWS profile, GCP zone, or Azure subscription")
	namespace := flag.String("namespace", "default", "The namespace to collect pod logs from")
	kubeconfigPath := flag.String("kubeconfigPath", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Path to the kubeconfig file")
	queryID := flag.String("queryid", "", "The queryID to filter pod logs from")
	podNames := flag.String("pods", "planner, executor, storage", "The specific pods to collect logs from (optional)")
	outputDir := flag.String("output", "./logs", "The output directory to save logs")
	flag.Parse()

	// Fetch kubeconfig for the specified cloud and cluster
	if err := utils.FetchKubeconfig(*cloudType, *clusterName, *authDetails, *kubeconfigPath); err != nil {
		fmt.Printf("Failed to fetch kubeconfig: %v\n", err)
		return
	}

	// Set up Kubernetes client
	clientset, err := utils.GetKubeClient(*kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to create Kubernetes client: %v\n", err)
		return
	}

	if *queryID != "" && *podNames != "" {
		if err := utils.SaveFilteredLogsForPods(clientset, *queryID, *namespace, *outputDir, *podNames); err != nil {
			fmt.Printf("Error collecting logs for pod %s: %v\n", *podNames, err)
		}
	} else {
		if err := utils.GetClusterLogs(clientset, *namespace, *outputDir); err != nil {
			fmt.Printf("Error collecting logs: %v\n", err)
		}
	}

}
