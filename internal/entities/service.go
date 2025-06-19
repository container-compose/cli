package entities

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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
func (service *Service) RunCommand(ctx context.Context) (*commands.RunCommand, error) {

	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return nil, err
		}
		service.Name = generated
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
func (service *Service) StartCommand(ctx context.Context) (*commands.StartCommand, error) {

	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return nil, err
		}
		service.Name = generated
	}

	cmd, err := commands.Start(service.Name)
	if err != nil {
		return nil, err
	}

	return cmd, nil
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
