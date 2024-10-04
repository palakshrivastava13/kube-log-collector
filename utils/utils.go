package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateLogFile creates a directory structure and file for saving pod logs
func CreateLogFile(directory, podName string) (*os.File, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// Define the path for the pod log file
	filePath := filepath.Join(directory, fmt.Sprintf("%s.log", podName))

	// Create a log file for the pod
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return file, nil
}
