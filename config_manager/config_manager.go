package config_manager

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/donovan-tuix/gman/config_paths"
)

// ConfigManager is a struct that holds the paths to the SSH config file and the SSH key file, as well as the file descriptor for the SSH config file.
type ConfigManager struct {
	Paths config_paths.ConfigPaths
	File  *os.File
}

var identityFileRegex = regexp.MustCompile(`^\s*IdentityFile\s+`)

// NewConfigManager creates a new ConfigManager struct and returns a pointer to it. It takes a newAccountName string as an argument, which is the name of the account that the user wants to use to authenticate with GitHub.
func NewConfigManager(newAccountName string) (*ConfigManager, error) {
	paths, err := config_paths.NewConfigPaths(newAccountName)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(paths.ConfigPath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &ConfigManager{
		Paths: *paths,
		File:  file,
	}, nil
}

// Update reads the SSH config file and replaces the IdentityFile line with the path to the SSH key file. It returns a slice of strings, which is the new SSH config file, and an error, if any.
func (cm *ConfigManager) Update() ([]string, error) {
	var newConfig []string
	scanner := bufio.NewScanner(cm.File)
	inGitHubHostSection := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Host github.com") {
			inGitHubHostSection = true
		} else if inGitHubHostSection && strings.HasPrefix(line, "Host ") {
			inGitHubHostSection = false
		}

		if inGitHubHostSection && identityFileRegex.MatchString(line) {
			line = fmt.Sprintf("	IdentityFile %s", cm.Paths.KeyPath)
		}

		newConfig = append(newConfig, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	cm.File.Truncate(0)
	cm.File.Seek(0, 0)
	return newConfig, nil
}

// WriteToFile writes the new SSH config file to the SSH config file. It takes a slice of strings as an argument, which is the new SSH config file. It returns an error, if any.
func (cm *ConfigManager) WriteToFile(lines []string) error {
	newConfigStr := strings.Join(lines, "\n")
	_, err := io.WriteString(cm.File, newConfigStr+"\n")
	return err
}

// Close closes the file descriptor for the SSH config file.
func (cm *ConfigManager) Close() {
	cm.File.Close()
}
