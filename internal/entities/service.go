package entities

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/big"

	"github.com/container-compose/cli/internal/commands"
	"github.com/goombaio/namegenerator"
	"gopkg.in/yaml.v3"
)

type Service struct {
	Image       string            `yaml:"image"`
	Name        string            `yaml:"container_name"`
	Ports       []string          `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
	Labels      map[string]string `yaml:"labels"`
	Volumes     []string          `yaml:"volumes"`
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

// RunCommand creates a command to run the service.
func (service *Service) RunCommand(ctx context.Context) (*commands.RunCommand, error) {

	if service.Name == "" {
		generated, err := service.GenerateName(ctx)
		if err != nil {
			return nil, err
		}
		service.Name = generated
	}

	cmd, err := commands.Run(service.Name)
	if err != nil {
		return nil, err
	}

	cmd.Image(service.Image)

	return cmd, nil
}
