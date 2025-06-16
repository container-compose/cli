package entities

import (
	"gopkg.in/yaml.v3"
)

type Compose struct {
	Version  string              `yaml:"version"`
	Services map[string]*Service `yaml:"services"`
}

func Parse(content []byte) (Compose, error) {
	config := Compose{}
	err := yaml.Unmarshal(content, &config)
	return config, err
}
