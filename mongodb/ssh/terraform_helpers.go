package ssh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/util"
)

// NewOpenPortCheckerFunc constructs a function based on the specified ssh.Client, which checks if the specified port is open
func NewOpenPortCheckerFunc(client *Client) func(port int) Result {
	return func(port int) Result {
		return client.RunCommand(fmt.Sprintf("(netstat -nl | grep -q %d) && echo open || echo closed", port))
	}
}

// IsPortOpen returns a StateRefreshFunc for determining if the specified port is open
func IsPortOpen(portChecker func(port int) Result, port int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result := portChecker(port)
		if result.IsError() {
			return nil, "", result
		}

		if result.Stdout == "open" {
			return port, "open", nil
		}

		if result.Stdout == "closed" {
			return nil, "closed", nil
		}

		if result.Stderr != "" {
			log.Printf("[DEBUG] Unexpected error: %s", result.Stderr)
		}

		return nil, "error", nil
	}
}

// WaitForOpenPort returns when the specified port is open, or with an error if the operation times out
func WaitForOpenPort(portChecker func(port int) Result, port int) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"closed"},
		Target:  []string{"open"},
		Refresh: IsPortOpen(portChecker, port),
		Timeout: util.PortAvailableTimeout,
	}

	log.Printf("[DEBUG] Waiting for port to be opened: %d", port)
	_, err := stateConf.WaitForState()

	return err
}

// NewServiceStatusChecker constructs a function based on the specified ssh.Client, which checks if the specified service is running
func NewServiceStatusChecker(client *Client) func(serviceName string) Result {
	return func(serviceName string) Result {
		return client.RunCommand(fmt.Sprintf("(ps -ef | grep -q %s && echo \"started\" ) || echo \"stopped\"", serviceName))
	}
}

// IsServiceRunning returns a StateRefreshFunc for determining if the specified service is running
func IsServiceRunning(serviceChecker func(serviceName string) Result, serviceName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result := serviceChecker(serviceName)
		if result.IsError() {
			return nil, "", result
		}

		if result.Stdout == "started" {
			return serviceName, "started", nil
		}

		if result.Stdout == "stopped" {
			return nil, "stopped", nil
		}

		if result.Stderr != "" {
			log.Printf("[DEBUG] Unexpected error: %s", result.Stderr)
		}

		return nil, "error", nil
	}
}

// WaitForService returns when the specified service is found to be running, or with an error if the operation times out
func WaitForService(serviceChecker func(serviceName string) Result, serviceName string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"stopped"},
		Target:  []string{"started"},
		Refresh: IsServiceRunning(serviceChecker, serviceName),
		Timeout: util.ServiceStartedTimeout,
	}

	log.Printf("[DEBUG] Waiting for service to start: %s", serviceName)
	_, err := stateConf.WaitForState()

	return err
}
