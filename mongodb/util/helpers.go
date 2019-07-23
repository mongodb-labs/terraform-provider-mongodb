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
// - if the passed contents are not an email, both the first/last names will be set to the specified value
//   and the email will be left blank
// if the email address does not contain a separator character, its user part will be set to both first/last name values
// if the email address has a dot or underscore separator character, the first name will be set to
//   the string up to the separator and the last name will be set to the character string starting after the separator
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

	// set the first name
	firstName = tentative[0]
	lastName = tentative[0]

	if len(tentative) > 1 {
		// if a last name was identified, set it
		lastName = tentative[1]
	}

	return
}
