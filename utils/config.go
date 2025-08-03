package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/copyleftdev/specgrade/core"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from YAML file or creates default config
func LoadConfig(configPath string) (*core.Config, error) {
	config := &core.Config{
		SpecVersion:   "3.1.0",
		OutputFormat:  "cli",
		FailThreshold: "B",
		SkipRules:     []string{},
	}

	// If no config path specified, look for specgrade.yaml in current directory
	if configPath == "" {
		if fileExists("specgrade.yaml") {
			configPath = "specgrade.yaml"
		} else if fileExists("specgrade.yml") {
			configPath = "specgrade.yml"
		}
	}

	// If still no config file, return default config
	if configPath == "" {
		return config, nil
	}

	// Check if config file exists
	if !fileExists(configPath) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// Read and parse config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	config.ConfigPath = configPath
	return config, nil
}

// MergeConfigWithFlags merges YAML config with CLI flags, giving precedence to CLI flags
func MergeConfigWithFlags(config *core.Config, flags *core.Config) *core.Config {
	merged := *config // Copy the config

	// CLI flags take precedence over config file
	if flags.SpecVersion != "" {
		merged.SpecVersion = flags.SpecVersion
	}
	if flags.InputDir != "" {
		merged.InputDir = flags.InputDir
	}
	if flags.OutputFormat != "" {
		merged.OutputFormat = flags.OutputFormat
	}
	if flags.FailThreshold != "" {
		merged.FailThreshold = flags.FailThreshold
	}
	if len(flags.SkipRules) > 0 {
		merged.SkipRules = flags.SkipRules
	}

	return &merged
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
