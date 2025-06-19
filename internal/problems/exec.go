package problems

import (
	"strings"
)

const (
	Generic = "000"
	Run     = "001"
	Stop    = "002"
	Inspect = "003"
	Start   = "004"
	Build   = "005"
)

var (
	ErrGeneric = New(Generic, "001", "An unknown error occurred")

	// run errors
	ErrNameCannotBeEmpty            = New(Run, "001", "Name cannot be empty")
	ErrContainerWithIDAlreadyExists = New(Run, "002", "A container with the same id already exists")

	// stop errors
	ErrIDCannotBeEmpty = New(Stop, "001", "ID cannot be empty")

	// inspect errors
	ErrContainerNotFound = New(Inspect, "001", "Container not found")

	// start errors
	ErrContainerAlreadyStarted = New(Start, "001", "Container is already started")

	// build errors
	ErrDockerfileNotFound   = New(Build, "001", "Dockerfile not found")
	ErrBuildContextNotFound = New(Build, "002", "Build context directory not found")
	ErrBuildFailed          = New(Build, "003", "Build failed")
)

var (
	// errorMap is a map of error messages to problems. The error messages are the ones that are returned
	// from the container engine.
	errorMap = map[string]Problem{
		"Error: exists: \"container with id ": ErrContainerWithIDAlreadyExists,
		"Error: not found: ":                  ErrContainerNotFound,
		"Error: already started: ":            ErrContainerAlreadyStarted,
		"Error: dockerfile not found":         ErrDockerfileNotFound,
		"Error: build context not found":      ErrBuildContextNotFound,
		"Error: build failed":                 ErrBuildFailed,
	}
)

func Convert(input string) Problem {
	for prefix, problem := range errorMap {
		if strings.HasPrefix(input, prefix) {
			return problem
		}
	}

	print(input)

	return ErrGeneric
}
