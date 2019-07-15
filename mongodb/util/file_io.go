package util

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ReadAllIntoTempFile reads all available data and saves it into a temporary file, using the specified pattern to identify it on the file system
// returns a file ready for reading
func ReadAllIntoTempFile(src io.Reader, pattern string) (file *os.File, err error) {
	file, err = ioutil.TempFile("", filepath.Clean(pattern)+"-*")
	if err != nil {
		return nil, err
	}
	defer SyncAndCloseFunc(file, &err)

	count, err := io.Copy(file, src)
	log.Printf("[DEBUG] copied %d bytes to %s", count, file.Name())
	if err != nil {
		return nil, err
	}

	return os.Open(file.Name())
}

// ReopenFile reopens the file specified by the passed file descriptor, for reading
func ReopenFile(file *os.File) (*os.File, error) {
	var err error
	filename := file.Name()
	if file, err = os.Open(filename); err != nil {
		return nil, fmt.Errorf("cannot reopen file %s: %v", filename, err)
	}
	_, _ = file.Seek(0, 0)

	return file, err
}

// SyncAndCloseFunc returns a function which can be passed to defer to sync and close the file, thus ensuring all data was saved on disk, or outlining any problems with the underlying file-system
// this implementation is loosely based on https://www.joeshaw.org/dont-defer-close-on-writable-files/ and https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling-overview.md
func SyncAndCloseFunc(file *os.File, topError *error) {
	log.Print("[DEBUG] Flushing all data to disk and closing the file")
	defer func() {
		if *topError != nil {
			// if the operation resulted in an error, remove the temporary file from disk
			log.Printf("[WARN] Deleting the file (%s), as an error was detected: %v", file.Name(), *topError)
			if err := os.Remove(file.Name()); err != nil {
				e := MergeErrors("failed to remove "+file.Name(), err, *topError)
				topError = &e
			}
		}
	}()

	if err := file.Sync(); err != nil {
		e := MergeErrors("failed to sync", err, *topError)
		topError = &e
		return
	}

	if err := file.Close(); err != nil {
		e := MergeErrors("failed to close", err, *topError)
		topError = &e
		return
	}
}

// BurnAfterReading ensure a consumed file is closed and deleted
func BurnAfterReading(file *os.File) {
	if err := file.Close(); err != nil {
		LogNonNilError(fmt.Errorf("failed to close %s: %v", file.Name(), err))
	}

	if err := os.Remove(file.Name()); err != nil {
		LogNonNilError(fmt.Errorf("failed to remove %s: %v", file.Name(), err))
	}
}

// MergeErrors constructs a new error based on a description string, given error, and a potentially previously existing error
// TODO(mihaibojin): https://godoc.org/gopkg.in/errgo.v2/fmt/errors
func MergeErrors(message string, newError error, previousError error) error {
	msg := fmt.Sprintf(message+": %v", newError)
	if previousError != nil {
		msg = msg + fmt.Sprintf(", previous err=%v", previousError)
	}

	return errors.New(msg)
}
