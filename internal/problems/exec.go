package problems

import (
	"strings"
)

const (
	Generic = "000"
	Run     = "001"
)

var (
	ErrGeneric = New(Generic, "001", "An unknown error occurred")

	// run errors
	ErrNameCannotBeEmpty            = New(Run, "001", "Name cannot be empty")
	ErrContainerWithIDAlreadyExists = New(Run, "002", "A container with the same id already exists")
)

var (
	// errorMap is a map of error messages to problems. The error messages are the ones that are returned
	// from the container engine.
	errorMap = map[string]Problem{
		"Error: exists: \"container with id ": ErrContainerWithIDAlreadyExists,
	}
)

func Convert(input string) Problem {
	for prefix, problem := range errorMap {
		if strings.HasPrefix(input, prefix) {
			return problem
		}
	}

	return ErrGeneric
}
