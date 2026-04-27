package state

import (
	"encoding/json"
	"fmt"
	"os"
)

// TerraformState represents the top-level structure of a Terraform state file.
type TerraformState struct {
	Version          int        `json:"version"`
	TerraformVersion string     `json:"terraform_version"`
	Resources        []Resource `json:"resources"`
}

// Resource represents a single resource block within the state file.
type Resource struct {
	Module    string     `json:"module,omitempty"`
	Mode      string     `json:"mode"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Provider  string     `json:"provider"`
	Instances []Instance `json:"instances"`
}

// Instance holds the actual attribute values for a resource instance.
type Instance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
}

// ParseStateFile reads and parses a Terraform state file from the given path.
func ParseStateFile(path string) (*TerraformState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file %q: %w", path, err)
	}

	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing state file %q: %w", path, err)
	}

	if state.Version < 4 {
		return nil, fmt.Errorf("unsupported state file version %d (minimum: 4)", state.Version)
	}

	return &state, nil
}

// ManagedResources returns only resources with mode "managed" (excludes data sources).
func (s *TerraformState) ManagedResources() []Resource {
	var managed []Resource
	for _, r := range s.Resources {
		if r.Mode == "managed" {
			managed = append(managed, r)
		}
	}
	return managed
}
