package config

import (
	"embed"
	"errors"
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

//go:embed config.yaml
var defaultConfig embed.FS

type Config struct {
	Name    string `yaml:"name"`
	Bucket  string `yaml:"bucket"`
	BaseDir string `yaml:"baseDir,omitempty"`

	Cli struct {
		Dev bool `yaml:"dev,omitempty"`
	} `yaml:"cli,omitempty"`

	Container struct {
		Image    string `yaml:"image"`
		Rootless bool   `yaml:"rootless,omitempty"`
	} `yaml:"container"`
}

func (c *Config) Read() error {
	file := simulatorConfigFile()

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		dir := filepath.Dir(file)
		if _, err := os.Stat(dir); err != nil {
			err = os.MkdirAll(dir, ownerReadWriteExecute)
			if err != nil {
				return errors.Join(errors.New("failed to create directory"), err)
			}
		}

		config, err := defaultConfig.ReadFile(FileName)
		if err != nil {
			return errors.Join(errors.New("failed to read config"), err)
		}

		err = os.WriteFile(file, config, ownerReadWrite)
		if err != nil {
			return errors.Join(errors.New("failed to write config"), err)
		}
	}

	config, err := os.ReadFile(file)
	if err != nil {
		return errors.Join(errors.New("failed to read config"), err)
	}

	err = yaml.Unmarshal(config, &c)
	if err != nil {
		return errors.Join(errors.New("failed to unmarshall config"), err)
	}

	return nil
}

func (c *Config) Write() error {
	config, err := yaml.Marshal(&c)
	if err != nil {
		return errors.Join(errors.New("failed to unmarshall config"), err)
	}

	if err := os.WriteFile(simulatorConfigFile(), config, ownerReadWrite); err != nil {
		return errors.Join(errors.New("failed to write config"), err)
	}

	return nil
}

func simulatorConfigFile() string {
	dir, ok := os.LookupEnv(Dir)
	if !ok {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		dir = filepath.Join(home, ".simulator")
	}

	return filepath.Join(dir, FileName)
}
