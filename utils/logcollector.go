package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPodLog collects logs from a particular pod and saves them in a file
func GetPodLog(clientset *kubernetes.Clientset, namespace, podName, outputDir string) error {
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

// GetClusterLogs GetPodLogs collects logs from all pods in a namespace and saves them in a directory
func GetClusterLogs(clientset *kubernetes.Clientset, namespace, outputDir string) error {
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
		if err := GetPodLog(clientset, namespace, pod.Name, namespaceDir); err != nil {
			fmt.Printf("Error collecting logs for pod %s: %v\n", pod.Name, err)
		}
	}

	return nil
}

// saveFilteredLogsForPod collects logs from a particular pod, filters them and saves them in a file
func saveFilteredLogsForPod(clientSet *kubernetes.Clientset, queryID, namespace, podName, outputDir string) error {
	// Get the logs of the specified pod
	logOptions := &v1.PodLogOptions{}
	req := clientSet.CoreV1().Pods(namespace).GetLogs(podName, logOptions)

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

	// Filter logs by queryID and write to file
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, queryID) {
			if _, err := file.WriteString(line + "\n"); err != nil {
				return fmt.Errorf("failed to write filtered log to file for pod %s: %v", podName, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading logs for pod %s: %v", podName, err)
	}

	fmt.Printf("Filtered logs for pod %s saved to %s\n", podName, file.Name())
	return nil
}

func SaveFilteredLogsForPods(clientSet *kubernetes.Clientset, queryID, namespace, outputDir, podNames string) error {
	combinedLogFilePath := filepath.Join(outputDir, "combined_filtered_logs.txt")
	combinedLogs := []string{}

	// Split the string by comma and trim spaces
	pods := strings.Split(podNames, ",")

	for _, podName := range pods {
		podName = strings.TrimSpace(podName)
		// Temporary log file for each pod
		tempFilePath := filepath.Join(outputDir, fmt.Sprintf("%s_filtered_logs.txt", podName))

		// Call the function to save filtered logs for each pod
		err := saveFilteredLogsForPod(clientSet, queryID, namespace, podName, tempFilePath)
		if err != nil {
			return fmt.Errorf("failed to save filtered logs for pod %s: %v", podName, err)
		}

		// Read filtered logs and append to combined logs
		tempFile, err := os.Open(tempFilePath)
		if err != nil {
			return fmt.Errorf("failed to open temp log file for pod %s: %v", podName, err)
		}
		defer tempFile.Close()

		scanner := bufio.NewScanner(tempFile)
		for scanner.Scan() {
			combinedLogs = append(combinedLogs, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading temp log file for pod %s: %v", podName, err)
		}
	}

	// Sort combined logs by timestamp
	sort.Slice(combinedLogs, func(i, j int) bool {
		timeI, errI := ExtractTimestamp(combinedLogs[i])
		timeJ, errJ := ExtractTimestamp(combinedLogs[j])
		if errI != nil || errJ != nil {
			return i < j // Default to original order if timestamp extraction fails
		}
		return timeI.Before(timeJ)
	})

	// Save sorted logs to the combined file
	combinedFile, err := os.Create(combinedLogFilePath)
	if err != nil {
		return fmt.Errorf("failed to create combined log file: %v", err)
	}
	defer combinedFile.Close()

	for _, line := range combinedLogs {
		if _, err := combinedFile.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to combined log file: %v", err)
		}
	}

	fmt.Printf("Combined and sorted logs saved to %s\n", combinedLogFilePath)
	return nil
}
