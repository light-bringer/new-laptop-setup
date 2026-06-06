package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Profiles map[string]ProfileAlias `yaml:"profiles"`
	Defaults Defaults                `yaml:"defaults"`
}

type ProfileAlias struct {
	AWSProfile  string `yaml:"aws_profile"`
	AccountID   string `yaml:"account_id"`
	Region      string `yaml:"region"`
	Description string `yaml:"description"`
}

type Defaults struct {
	Region string `yaml:"region"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return &Config{}, nil
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.Profiles == nil {
		cfg.Profiles = map[string]ProfileAlias{}
	}

	return cfg, nil
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".awsx.yaml")
}
