package util

import (
	"time"
)

const (
	// DefaultTimeout default timeout for most operations
	DefaultTimeout = 10 * time.Minute

	// DownloadTimeout is the configured download timeout after which the run is aborted
	DownloadTimeout = 30 * time.Minute

	// PortAvailableTimeout timeout used while waiting for a port to become available
	PortAvailableTimeout = 30 * time.Minute

	// LongCreationTimeout timeout to use for longer operations
	LongCreationTimeout = 60 * time.Minute
)
