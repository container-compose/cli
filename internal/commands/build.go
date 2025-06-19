package commands

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/container-compose/cli/internal/problems"
)

type BuildCommand struct {
	Context    string
	Dockerfile string
	Tag        string
	BuildArgs  map[string]string
	Labels     map[string]string
	Target     string
	NoCache    bool
	Pull       bool
	Quiet      bool
	CPUs       int
	Memory     string
	Arch       string
	OS         string
	Progress   string
}

func Build(context string) (*BuildCommand, error) {
	if context == "" {
		context = "."
	}

	return &BuildCommand{
		Context:  context,
		CPUs:     2,
		Memory:   "2048MB",
		Arch:     "arm64",
		OS:       "linux",
		Progress: "auto",
	}, nil
}

// SetDockerfile sets the path to the Dockerfile
func (c *BuildCommand) SetDockerfile(dockerfile string) *BuildCommand {
	c.Dockerfile = dockerfile
	return c
}

// SetTag sets the tag for the built image
func (c *BuildCommand) SetTag(tag string) *BuildCommand {
	c.Tag = tag
	return c
}

// SetBuildArgs sets build-time variables
func (c *BuildCommand) SetBuildArgs(args map[string]string) *BuildCommand {
	c.BuildArgs = args
	return c
}

// SetLabels sets labels for the built image
func (c *BuildCommand) SetLabels(labels map[string]string) *BuildCommand {
	c.Labels = labels
	return c
}

// SetTarget sets the target build stage
func (c *BuildCommand) SetTarget(target string) *BuildCommand {
	c.Target = target
	return c
}

// SetNoCache disables cache usage
func (c *BuildCommand) SetNoCache(noCache bool) *BuildCommand {
	c.NoCache = noCache
	return c
}

// SetPull always attempts to pull newer versions of base images
func (c *BuildCommand) SetPull(pull bool) *BuildCommand {
	c.Pull = pull
	return c
}

// SetQuiet suppresses build output
func (c *BuildCommand) SetQuiet(quiet bool) *BuildCommand {
	c.Quiet = quiet
	return c
}

// SetCPUs sets the number of CPUs to allocate
func (c *BuildCommand) SetCPUs(cpus int) *BuildCommand {
	c.CPUs = cpus
	return c
}

// SetMemory sets the memory limit
func (c *BuildCommand) SetMemory(memory string) *BuildCommand {
	c.Memory = memory
	return c
}

// SetArch sets the build architecture
func (c *BuildCommand) SetArch(arch string) *BuildCommand {
	c.Arch = arch
	return c
}

// SetOS sets the build OS
func (c *BuildCommand) SetOS(os string) *BuildCommand {
	c.OS = os
	return c
}

// SetProgress sets the progress type
func (c *BuildCommand) SetProgress(progress string) *BuildCommand {
	c.Progress = progress
	return c
}

// Exec executes the build command
func (c *BuildCommand) Exec(ctx context.Context) error {
	args := []string{
		"build",
	}

	// Add CPU allocation
	if c.CPUs > 0 {
		args = append(args, "--cpus", strconv.Itoa(c.CPUs))
	}

	// Add memory limit
	if c.Memory != "" {
		args = append(args, "--memory", c.Memory)
	}

	// Add build args
	for key, value := range c.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	// Add Dockerfile path
	if c.Dockerfile != "" {
		args = append(args, "--file", c.Dockerfile)
	}

	// Add labels
	for key, value := range c.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", key, value))
	}

	// Add no-cache flag
	if c.NoCache {
		args = append(args, "--no-cache")
	}

	// Add architecture
	if c.Arch != "" {
		args = append(args, "--arch", c.Arch)
	}

	// Add OS
	if c.OS != "" {
		args = append(args, "--os", c.OS)
	}

	// Add progress type
	if c.Progress != "" {
		args = append(args, "--progress", c.Progress)
	}

	// Add tag
	if c.Tag != "" {
		args = append(args, "--tag", c.Tag)
	}

	// Add target
	if c.Target != "" {
		args = append(args, "--target", c.Target)
	}

	// Add quiet flag
	if c.Quiet {
		args = append(args, "--quiet")
	}

	// Add context directory
	args = append(args, c.Context)

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
