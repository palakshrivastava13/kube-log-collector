# Kube Log Collector

`Kube Log Collector` is a Go-based tool designed to assist testers by automatically collecting Kubernetes pod logs during test failures. It can be integrated as a plugin into test suites to generate and save log files whenever a test case fails. This helps in faster troubleshooting by capturing the relevant pod logs in real time.

## Features

- **Test Integration**: Can be used as a plugin within test frameworks to collect logs automatically when tests fail.
- **Collect logs for all pods in a namespace**: Fetch logs from all pods within a specified namespace when triggered by a test case.
- **Collect logs for a specific pod**: Fetch logs from a specific pod based on test case context.
- **Customizable output directory**: Save logs in a user-defined directory structure, making it easy to access logs related to specific test cases.
- **Directory structure**: Organizes pod logs in a structured directory format for easy navigation and debugging.

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/kube-log-collector.git
   cd kube-log-collector
2. **Install dependencies**:

   Make sure you have Go installed. If not, install Go.
   
   In the project root, run:

   ```bash
   go mod tidy
   ```
   This will download and install the necessary dependencies.

## Usage
To integrate kube-log-collector as a plugin in your test framework:

1. **Set up the plugin**:

   Import the kube-log-collector/pkg/logs package into your test project.
   ```bash
   import "kube-log-collector/pkg/logs"

2. **Trigger log collection on test failure**:

   Inside your test framework’s teardown function or failure handler, call the appropriate log collection function, passing in the relevant test context (e.g., namespace, pod name).

## Example in a Test Framework
Here’s an example of how you can use it within a test suite (such as in Go):

```bash
func TestSomeFunction(t *testing.T) {
    defer func() {
        if t.Failed() {
            // Collect logs on test failure
            err := logs.GetPodLogs("/path/to/kubeconfig", "test-namespace", "./logs")
            if err != nil {
                t.Errorf("Failed to collect logs: %v", err)
            }
        }
    }()
    
    // Your test code here
}
```
Alternatively, to collect logs for a specific pod:

```bash
func TestSpecificPod(t *testing.T) {
    defer func() {
        if t.Failed() {
            // Collect logs for a specific pod on test failure
            err := logs.GetPodLog("/path/to/kubeconfig", "test-namespace", "test-pod", "./logs")
            if err != nil {
                t.Errorf("Failed to collect logs: %v", err)
            }
        }
    }()
    
    // Your test code here
}

```

# Command-line Usage for Manual Testing
You can also run the tool manually via the command-line to collect logs outside of automated tests.

1. **Collect logs for all pods in a namespace**:
   ```bash
   go run cmd/kube-log-collector/main.go \
   --kubeconfig=/path/to/kubeconfig \
   --namespace=my-namespace \
   --output=/path/to/output-directory
   
2. **Collect logs for a specific pod**:
   ```bash
   go run cmd/kube-log-collector/main.go \
   --kubeconfig=/path/to/kubeconfig \
   --namespace=my-namespace \
   --pod=my-pod-name \
   --output=/path/to/output-directory