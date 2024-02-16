package config_paths

import (
	"os/user"
	"path/filepath"
)

type ConfigPaths struct {
	ConfigPath string
	KeyPath    string
}

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
