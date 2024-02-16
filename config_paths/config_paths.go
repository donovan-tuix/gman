package config_paths

import (
	"os/user"
	"path/filepath"
)

// ConfigPaths is a struct that holds the paths to the SSH config file and the SSH key file.
type ConfigPaths struct {
	ConfigPath string
	KeyPath    string
}

// NewConfigPaths creates a new ConfigPaths struct and returns a pointer to it. It takes a newAccountName string as an argument, which is the name of the account that the user wants to use to authenticate with GitHub.
func NewConfigPaths(newAccountName string) (*ConfigPaths, error) {
	homeDir, err := user.Current()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir.HomeDir, ".ssh/config")
	keyPath := filepath.Join(homeDir.HomeDir, ".ssh", newAccountName)

	return &ConfigPaths{
		ConfigPath: configPath,
		KeyPath:    keyPath,
	}, nil
}
