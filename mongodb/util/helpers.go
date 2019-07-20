package util

import (
	"log"
	"regexp"
	"strings"
)

var reValidEmail *regexp.Regexp

func init() {
	reValidEmail = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,32}$`)
}

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

// IsValidEmailAddress returns true if the passed string is a valid email address
func IsValidEmailAddress(email string) bool {
	return reValidEmail.MatchString(email)
}

// TryExtractFirstLastNameAndEmail attempts to extract first, last name and email address
func TryExtractFirstLastNameAndEmail(contents string) (firstName string, lastName string, emailAddress string) {
	if !IsValidEmailAddress(contents) {
		// if the data string is not an email, set first/last name to the passed contents
		firstName = contents
		lastName = contents
		emailAddress = ""
		return
	}
	emailAddress = contents

	part := strings.Split(contents, "@")[0]
	tentative := strings.SplitN(part, ".", 2)
	if len(tentative) < 2 {
		// if the contents does not contain a dot, attempt to split by underscore
		tentative = strings.SplitN(part, "_", 2)
	}

	if len(tentative) < 2 {
		// could not find 2 parts, set both as the same
		firstName = tentative[0]
		lastName = tentative[0]
		emailAddress = contents
		return
	}

	firstName = tentative[0]
	lastName = tentative[1]
	return
}
