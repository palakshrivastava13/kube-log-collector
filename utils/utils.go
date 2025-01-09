package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// Helper function to extract timestamp from a log line
func ExtractTimestamp(logLine string) (time.Time, error) {
	// Assume timestamp is the first field in the log line, e.g., "2025-01-09T10:45:00Z ..."
	parts := strings.Fields(logLine)
	if len(parts) == 0 {
		return time.Time{}, fmt.Errorf("invalid log line")
	}
	return time.Parse(time.RFC3339, parts[0])
}
