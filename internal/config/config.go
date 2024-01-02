package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

const (
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	XDGConfigHomeLinuxEnv  = "XDG_CONFIG_HOME"
	LocalAppDataWindowsEnv = "LOCALAPPDATA"
	Dir                    = "SIMULATOR_DIR"
	FileName               = "config.yaml"
	ownerReadWrite         = 0600
	ownerReadWriteExecute  = 0700
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
		err = os.MkdirAll(dir, ownerReadWriteExecute)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
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
	dir, err := SimulatorDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "admin"), nil
}

func (c *Config) PlayerBundleDir() (string, error) {
	dir, err := SimulatorDir()
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

func SimulatorDir() (string, error) {
	// User provided config has precedence
	if dir, ok := os.LookupEnv(Dir); ok {
		return dir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user's home directory: %w", err)
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Preferences", "io.controlplane.simulator"), nil
	case "linux":
		if dir, ok := os.LookupEnv(XDGConfigHomeLinuxEnv); ok {
			return filepath.Join(dir, "simulator"), nil
		}

		// Fallback to default XDG dir
		return filepath.Join(homeDir, ".config", "simulator"), nil
	case "windows":
		if dir, ok := os.LookupEnv(LocalAppDataWindowsEnv); ok {
			return filepath.Join(dir, "simulator"), nil
		}

		// Fallback to default local app data dir
		return filepath.Join(homeDir, "AppData", "Local", "simulator"), nil
	default:
		return "", fmt.Errorf("operating system not support: %s", runtime.GOOS)
	}
}

func simulatorConfigFile() (string, error) {
	dir, err := SimulatorDir()
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
