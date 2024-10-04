package main

import (
	"flag"
	"fmt"
	logs "kube-log-collector/utils"
	"os"
	"path/filepath"
)

func main() {
	// Define command-line flags for kubeconfig file, namespace, pod name, and output directory
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Path to the kubeconfig file")
	namespace := flag.String("namespace", "default", "The namespace to collect pod logs from")
	podName := flag.String("pod", "", "The specific pod to collect logs from (optional)")
	outputDir := flag.String("output", "./logs", "The output directory to save logs")
	flag.Parse()

	// If podName is provided, collect logs for the specific pod, otherwise collect logs for all pods in the namespace
	if *podName != "" {
		if err := logs.GetPodLog(*kubeconfig, *namespace, *podName, *outputDir); err != nil {
			fmt.Printf("Error collecting logs for pod %s: %v\n", *podName, err)
		}
	} else {
		if err := logs.GetPodLogs(*kubeconfig, *namespace, *outputDir); err != nil {
			fmt.Printf("Error collecting logs: %v\n", err)
		}
	}
}
