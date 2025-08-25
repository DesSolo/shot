package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Screenshotter ...
type Screenshotter struct {
	Enabled    bool   `yaml:"enabled"`
	Host       string `yaml:"host"`
	UserAgent  string `yaml:"user_agent"`
	Resolution *struct {
		Width  int64 `yaml:"width"`
		Height int64 `yaml:"height"`
	} `yaml:"resolution"`
	Cookies []string          `yaml:"cookies"`
	Headers map[string]string `yaml:"headers"`
}

// Config ...
type Config struct {
	Logger struct {
		Level int `yaml:"level"`
	} `yaml:"logger"`
	Screenshotter struct {
		Default *Screenshotter           `yaml:"default"`
		Sites   map[string]Screenshotter `yaml:"sites"`
	} `yaml:"screenshotter"`
}

// NewConfigFromFile ...
func NewConfigFromFile(file string) (*Config, error) {
	data, err := os.ReadFile(file) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return &config, nil
}
