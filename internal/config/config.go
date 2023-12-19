package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	Dir                   = "SIMULATOR_DIR"
	FileName              = "config.yaml"
	ownerReadWrite        = 0600
	ownerReadWriteExecute = 0700
)

type Config struct {
	Name      string `yaml:"name"`
	Bucket    string `yaml:"bucket"`
	BaseDir   string `yaml:"baseDir,omitempty"`
	Cli       `yaml:"cli,omitempty"`
	Container `yaml:"container"`
}

type Cli struct {
	Dev bool `yaml:"dev,omitempty"`
}

type Container struct {
	Image    string `yaml:"image"`
	Rootless bool   `yaml:"rootless,omitempty"`
}

func (c *Config) Read() error {
	file, err := simulatorConfigFile()
	if err != nil {
		return err
	}

	if _, err = os.Stat(file); errors.Is(err, os.ErrNotExist) {
		dir := filepath.Dir(file)
		if _, err = os.Stat(dir); err != nil {
			err = os.MkdirAll(dir, ownerReadWriteExecute)
			if err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}
		}

		config := defaultConfig()
		if err = config.Write(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		return nil
	}

	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return fmt.Errorf("failed to decode config to bytes: %w", err)
	}

	return nil
}

func (c *Config) Write() error {
	config, err := yaml.Marshal(&c)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	file, err := simulatorConfigFile()
	if err != nil {
		return err
	}

	if err = os.WriteFile(file, config, ownerReadWrite); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (c *Config) AdminBundleDir() (string, error) {
	dir, err := simulatorDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "admin"), nil
}

func (c *Config) PlayerBundleDir() (string, error) {
	dir, err := simulatorDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "player"), nil
}

func (c *Config) ContainerUser() string {
	if c.Rootless {
		return "root"
	}

	return "ubuntu"
}

func simulatorDir() (string, error) {
	dir, ok := os.LookupEnv(Dir)
	if !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to determine user's home directory: %w", err)
		}

		return filepath.Join(home, ".simulator"), nil
	}

	return dir, nil
}

func simulatorConfigFile() (string, error) {
	dir, err := simulatorDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, FileName), nil
}

func defaultConfig() Config {
	return Config{
		Name: "simulator",
		Container: Container{
			Image: "controlplane/simulator:latest",
		},
	}
}
