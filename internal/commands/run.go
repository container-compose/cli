package commands

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/container-compose/cli/internal/problems"
)

type RunCommand struct {
	Name   string
	Attach bool
	// Interactive    bool
	// Debug          bool
	// Version        bool
	ContainerImage string
}

func (c *RunCommand) Image(image string) *RunCommand {
	c.ContainerImage = image
	return c
}

func Run(name string) (*RunCommand, error) {
	if name == "" {
		return nil, problems.ErrNameCannotBeEmpty
	}

	return &RunCommand{
		Name: name,
	}, nil
}

// Exec executes the run command
func (c *RunCommand) Exec(ctx context.Context) error {

	args := []string{
		"run",
		"--name", c.Name,
	}

	if !c.Attach {
		args = append(args, "--detach")
	}

	args = append(args, c.ContainerImage)
	cmd := exec.Command("container", args...)

	// create io writers to capture the exec output
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return problems.Convert(stderr.String())
	}

	return nil
}
