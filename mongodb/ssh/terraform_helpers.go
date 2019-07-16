package ssh

import (
	"fmt"
	"log"

	"github.com/mongodb-labs/terraform-provider-mongodb/mongodb/util"

	"github.com/hashicorp/terraform/helper/resource"
)

// NewOpenPortCheckerFunc constructs a function based on the specified ssh.Client, which checks if the specified port is open
func NewOpenPortCheckerFunc(client *Client) func(port int) Result {
	return func(port int) Result {
		return client.RunCommand(fmt.Sprintf("(ss -tln | grep -q %d) && echo open || echo closed", port))
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
