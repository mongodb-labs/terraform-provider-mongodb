package util

import "log"

// PanicOnError panics when the action is run and results in an error
func PanicOnError(action func() error) {
	PanicOnNonNilErr(action())
}

// PanicOnNonNilErr convenience wrapper which panics when the passed error is not nil
func PanicOnNonNilErr(err error) {
	if err != nil {
		log.Fatalf("err=%v", err)
	}
}

// LogError logs any errors returned by the action
func LogError(action func() error) {
	LogNonNilError(action())
}

// LogNonNilError logs the error if not nil
func LogNonNilError(err error) {
	if err != nil {
		log.Printf("Unexpected error: %v", err)
	}
}
