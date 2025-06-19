package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"

	"github.com/container-compose/cli/internal/problems"
)

type InspectCommand struct {
	ID string
}

// InspectResult represents the JSON output from the container inspect command
type InspectResult struct {
	Configuration Configuration `json:"configuration"`
	Networks      []interface{} `json:"networks"`
	Status        string        `json:"status"`
}

type Configuration struct {
	DNS            DNS                    `json:"dns"`
	Networks       []string               `json:"networks"`
	Labels         map[string]string      `json:"labels"`
	Image          Image                  `json:"image"`
	Platform       Platform               `json:"platform"`
	Resources      Resources              `json:"resources"`
	InitProcess    InitProcess            `json:"initProcess"`
	Hostname       string                 `json:"hostname"`
	Mounts         []interface{}          `json:"mounts"`
	Rosetta        bool                   `json:"rosetta"`
	RuntimeHandler string                 `json:"runtimeHandler"`
	ID             string                 `json:"id"`
	Sysctls        map[string]interface{} `json:"sysctls"`
}

type DNS struct {
	Nameservers   []string      `json:"nameservers"`
	Domain        string        `json:"domain"`
	SearchDomains []interface{} `json:"searchDomains"`
	Options       []interface{} `json:"options"`
}

type Image struct {
	Descriptor Descriptor `json:"descriptor"`
	Reference  string     `json:"reference"`
}

type Descriptor struct {
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	MediaType string `json:"mediaType"`
}

type Platform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

type Resources struct {
	MemoryInBytes int64 `json:"memoryInBytes"`
	CPUs          int   `json:"cpus"`
}

type InitProcess struct {
	Arguments          []string      `json:"arguments"`
	Executable         string        `json:"executable"`
	Environment        []string      `json:"environment"`
	User               User          `json:"user"`
	WorkingDirectory   string        `json:"workingDirectory"`
	SupplementalGroups []interface{} `json:"supplementalGroups"`
	Terminal           bool          `json:"terminal"`
	Rlimits            []interface{} `json:"rlimits"`
}

type User struct {
	ID UserID `json:"id"`
}

type UserID struct {
	UID int `json:"uid"`
	GID int `json:"gid"`
}

func Inspect(id string) (*InspectCommand, error) {
	if id == "" {
		return nil, problems.ErrIDCannotBeEmpty
	}

	return &InspectCommand{
		ID: id,
	}, nil
}

// Exec executes the inspect command and returns the parsed result
func (c *InspectCommand) Exec(ctx context.Context) ([]InspectResult, error) {
	args := []string{
		"inspect",
		c.ID,
	}

	cmd := exec.Command("container", args...)

	// create io writers to capture the exec output
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return nil, problems.Convert(stderr.String())
	}

	// Parse the JSON output
	var results []InspectResult
	if err := json.Unmarshal(stdout.Bytes(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
