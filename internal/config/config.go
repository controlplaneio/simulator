package config

import (
	"embed"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	Dir      = "SIMULATOR_DIR"
	FileName = "config.yaml"
)

//go:embed config.yaml
var defaultConfig embed.FS

type Config struct {
	Name    string `yaml:"name"`
	Bucket  string `yaml:"bucket"`
	BaseDir string `yaml:"baseDir"`

	Cli struct {
		Dev bool `yaml:"dev,omitempty"`
	} `yaml:"cli,omitempty"`

	Container struct {
		Image string `yaml:"image"`
	} `yaml:"container"`
}

func (c *Config) Read() error {
	file := file()

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		dir := filepath.Dir(file)
		if _, err := os.Stat(dir); err != nil {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}

		config, err := defaultConfig.ReadFile(FileName)
		if err != nil {
			return err
		}

		err = os.WriteFile(file, config, 0644)
		if err != nil {
			return err
		}
	}

	config, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(config, &c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Write() error {
	config, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	return os.WriteFile(file(), config, 0644)
}

func file() string {
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
