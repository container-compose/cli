package commands

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/container-compose/cli/internal/problems"
)

type StopCommand struct {
	ID string
}

func Stop(id string) (*StopCommand, error) {
	if id == "" {
		return nil, problems.ErrIDCannotBeEmpty
	}

	return &StopCommand{
		ID: id,
	}, nil
}

// Exec executes the stop command
func (c *StopCommand) Exec(ctx context.Context) error {

	args := []string{
		"stop",
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
		return problems.Convert(stderr.String())
	}

	return nil
}
