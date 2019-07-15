package ssh

import (
	"fmt"
	"log"
	"strings"
)

// Result represents the result of a SSH command
type Result struct {
	Cmd    string
	Stdout string
	Stderr string
	Err    error
}

// Error implementation of the error interface
func (r Result) Error() string {
	return fmt.Sprintf("ssh.%s: unexpected error=%s", r.Cmd, r.GetDebugInfo())
}

// IsError true if the Result contains an error
func (r Result) IsError() bool {
	return r.Err != nil
}

// GetDebugInfo returns debugging information to be displayed where needed
func (r *Result) GetDebugInfo() string {
	var sb strings.Builder
	if r.Cmd != "" {
		_, _ = fmt.Fprintf(&sb, " command: '%s'", r.Cmd)
	}
	_, _ = fmt.Fprintf(&sb, "\nIsError: %v", r.IsError())

	if r.Stdout != "" {
		_, _ = fmt.Fprintf(&sb, "\nstdout: %s", r.Stdout)
	}

	if r.Stderr != "" {
		_, _ = fmt.Fprintf(&sb, "\nstderr: %s", r.Stderr)
	}

	return sb.String()
}

// Debug print debugging information and return the same object
func (r *Result) Debug() {
	log.Println("[DEBUG] Executed" + r.GetDebugInfo())
}

// PanicOnError panic if a result has completed with an error (useful in situations where errors are not recoverable,
// but we want to allow error handling at an upper layer
func PanicOnError(r Result) {
	if r.IsError() {
		log.Fatalf("unexpected error for %s; err=%v", r.GetDebugInfo(), r.Err)
	}
}
