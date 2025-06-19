package entities

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/container-compose/cli/internal/commands"
	"github.com/goombaio/namegenerator"
	"gopkg.in/yaml.v3"
)

type Service struct {
	Image                string            `yaml:"image"`
	Name                 string            `yaml:"container_name"`
	Ports                []string          `yaml:"ports"`
	EnvironmentVariables map[string]string `yaml:"environment"`
	Labels               map[string]string `yaml:"labels"`
	Volumes              []string          `yaml:"volumes"`
	Build                *Build            `yaml:"build,omitempty"`
}

type Build struct {
	Context    string            `yaml:"context,omitempty"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
	Labels     map[string]string `yaml:"labels,omitempty"`
	Target     string            `yaml:"target,omitempty"`
	Network    string            `yaml:"network,omitempty"`
	NoCache    bool              `yaml:"no_cache,omitempty"`
	Pull       bool              `yaml:"pull,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for Build which
// handles both string format (build: "./path") and object format
func (b *Build) UnmarshalYAML(value *yaml.Node) error {
	// If it's a string, treat it as the context path
	if value.Kind == yaml.ScalarNode {
		b.Context = value.Value
		return nil
	}

	// If it's a mapping, unmarshal normally
	if value.Kind == yaml.MappingNode {
		type buildAlias Build
		aux := (*buildAlias)(b)
		return value.Decode(aux)
	}

	return fmt.Errorf("build must be either a string or an object")
}

// GenerateName generates a name for the service. It does this deterministically based on the
// configuration of the service. This allows for consistent naming across multiple runs of the
// application and for commands which require stop / restart to function correctly.
//
// It is limited between runs when the service configuration changes as the generated name
// will be different.
func (service *Service) GenerateName(ctx context.Context) (string, error) {

	// marshall the service to a string
	data, err := yaml.Marshal(service)
	if err != nil {
		return "", err
	}

	// create a seed based on the hash of the service configuration
	seed := big.NewInt(0)
	hash := md5.New()
	hash.Write([]byte(data))
	hexstr := hex.EncodeToString(hash.Sum(nil))
	seed.SetString(hexstr, 16)

	// create a name generator based on the seed
	nameGenerator := namegenerator.NewNameGenerator(seed.Int64())
	return nameGenerator.Generate(), nil
}

// Exists checks if the service exists.
func (service *Service) Exists(ctx context.Context) (bool, error) {
	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return false, err
		}
		service.Name = generated
	}

	cmd, err := commands.Inspect(service.Name)
	if err != nil {
		return false, err
	}

	results, err := cmd.Exec(ctx)
	if err != nil {
		return false, nil // If inspect fails, assume the container doesn't exist
	}

	return len(results) > 0, nil
}

// IsRunning checks if the service is running.
func (service *Service) IsRunning(ctx context.Context) (bool, error) {
	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return false, err
		}
		service.Name = generated
	}

	cmd, err := commands.Inspect(service.Name)
	if err != nil {
		return false, err
	}

	results, err := cmd.Exec(ctx)
	if err != nil {
		return false, nil // If inspect fails, assume container is not running
	}

	if len(results) == 0 {
		return false, nil
	}

	return results[0].Status == "running", nil
}

// RunCommand creates a command to run the service.
// If the service has build configuration, it will build the image first.
func (service *Service) RunCommand(ctx context.Context) (*commands.RunCommand, error) {

	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return nil, err
		}
		service.Name = generated
	}

	// Check if we need to build the image first
	if service.Build != nil {
		buildCmd, err := service.BuildCommand(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create build command: %w", err)
		}

		// Execute the build command
		if err := buildCmd.Exec(ctx); err != nil {
			return nil, fmt.Errorf("failed to build image: %w", err)
		}

		// If no image was specified, use the built image tag
		if service.Image == "" {
			if service.Build.Context != "" {
				// Use the service name as the image tag
				service.Image = service.Name
			}
		}
	}

	cmd, err := commands.Run(service.Name, service.EnvironmentVariables, service.Labels)
	if err != nil {
		return nil, err
	}

	cmd.Image(service.Image)

	return cmd, nil
}

// StopCommand creates a command to stop the service.
func (service *Service) StopCommand(ctx context.Context) (*commands.StopCommand, error) {

	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return nil, err
		}
		service.Name = generated
	}

	cmd, err := commands.Stop(service.Name)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// StartCommand creates a command to start the service.
// If the service has build configuration and the image doesn't exist, it will build first.
func (service *Service) StartCommand(ctx context.Context) (*commands.StartCommand, error) {

	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return nil, err
		}
		service.Name = generated
	}

	// Check if we need to build the image first
	if service.Build != nil {
		// If we have a build configuration, ensure the image is built
		exists, err := service.Exists(ctx)
		if err != nil {
			return nil, err
		}

		if !exists {
			// Container doesn't exist, we need to build and run instead of start
			return nil, fmt.Errorf("container does not exist, use RunCommand to build and run")
		}
	}

	cmd, err := commands.Start(service.Name)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// BuildCommand creates a command to build the service image.
func (service *Service) BuildCommand(ctx context.Context) (*commands.BuildCommand, error) {
	if service.Build == nil {
		return nil, fmt.Errorf("no build configuration found for service")
	}

	// Use the build context, default to current directory
	context := service.Build.Context
	if context == "" {
		context = "."
	}

	cmd, err := commands.Build(context)
	if err != nil {
		return nil, err
	}

	// Set Dockerfile path if specified
	if service.Build.Dockerfile != "" {
		cmd.SetDockerfile(service.Build.Dockerfile)
	}

	// Set build args
	if service.Build.Args != nil {
		cmd.SetBuildArgs(service.Build.Args)
	}

	// Set labels (merge service labels with build labels)
	allLabels := make(map[string]string)
	for k, v := range service.Labels {
		allLabels[k] = v
	}
	for k, v := range service.Build.Labels {
		allLabels[k] = v
	}
	if len(allLabels) > 0 {
		cmd.SetLabels(allLabels)
	}

	// Set target
	if service.Build.Target != "" {
		cmd.SetTarget(service.Build.Target)
	}

	// Set no-cache
	if service.Build.NoCache {
		cmd.SetNoCache(true)
	}

	// Set pull
	if service.Build.Pull {
		cmd.SetPull(true)
	}

	// Set tag - use the service image name if specified, otherwise use service name
	tag := service.Image
	if tag == "" {
		if service.Name == "" {
			generated, err := service.GenerateName(ctx)
			if err != nil {
				return nil, err
			}
			tag = generated
		} else {
			tag = service.Name
		}
	}
	cmd.SetTag(tag)

	return cmd, nil
}

// NeedsBuild checks if the service needs to be built (has build config but no image)
func (service *Service) NeedsBuild() bool {
	return service.Build != nil && service.Image == ""
}

// InspectCommand creates a command to inspect the service.
func (service *Service) InspectCommand(ctx context.Context) (*commands.InspectCommand, error) {

	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return nil, err
		}
		service.Name = generated
	}

	cmd, err := commands.Inspect(service.Name)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// FromInspectResult creates a Service from an InspectResult. This allows us to convert
// the actual state of a running container back into our desired state representation.
func FromInspectResult(result commands.InspectResult) *Service {
	// Convert environment variables from slice to map
	envVars := make(map[string]string)
	for _, env := range result.Configuration.InitProcess.Environment {
		// Parse environment variables in the format KEY=VALUE
		if len(env) > 0 {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				envVars[parts[0]] = parts[1]
			}
		}
	}

	// Extract ports from labels or other configuration if available
	// This is a simplified implementation as the inspect format doesn't directly show port mappings
	var ports []string

	// Extract volumes from mounts if available
	var volumes []string

	return &Service{
		Image:                result.Configuration.Image.Reference,
		Name:                 result.Configuration.ID,
		Ports:                ports,
		EnvironmentVariables: envVars,
		Labels:               result.Configuration.Labels,
		Volumes:              volumes,
	}
}
