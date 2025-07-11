package commands

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/container-compose/cli/internal/problems"
)

type StartCommand struct {
	ID string
}

func Start(id string) (*StartCommand, error) {
	if id == "" {
		return nil, problems.ErrIDCannotBeEmpty
	}

	return &StartCommand{
		ID: id,
	}, nil
}

// Exec executes the start command
func (c *StartCommand) Exec(ctx context.Context) error {

	args := []string{
		"start",
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
