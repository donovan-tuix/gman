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

type ConfigManager struct {
	Paths config_paths.ConfigPaths
	File  *os.File
}

var identityFileRegex = regexp.MustCompile(`^\s*IdentityFile\s+`)

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

func (cm *ConfigManager) WriteToFile(lines []string) error {
	newConfigStr := strings.Join(lines, "\n")
	_, err := io.WriteString(cm.File, newConfigStr+"\n")
	return err
}

func (cm *ConfigManager) Close() {
	cm.File.Close()
}
