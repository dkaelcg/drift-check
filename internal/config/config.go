package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level drift-check configuration.
type Config struct {
	StateFile  string            `yaml:"state_file"`
	Region     string            `yaml:"region"`
	OutputFmt  string            `yaml:"output_format"`
	Filters    []string          `yaml:"filters"`
	Ignore     []string          `yaml:"ignore"`
	ExtraVars  map[string]string `yaml:"extra_vars"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		StateFile: "terraform.tfstate",
		Region:    "us-east-1",
		OutputFmt: "text",
	}
}

// Load reads a YAML config file from path and merges it over defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate performs basic sanity checks on the loaded configuration.
func (c *Config) validate() error {
	validFormats := map[string]bool{"text": true, "json": true}
	if !validFormats[c.OutputFmt] {
		return fmt.Errorf("unsupported output_format %q: must be one of text, json", c.OutputFmt)
	}
	if c.StateFile == "" {
		return fmt.Errorf("state_file must not be empty")
	}
	return nil
}
