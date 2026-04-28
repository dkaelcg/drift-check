package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "drift.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("/nonexistent/drift.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StateFile != "terraform.tfstate" {
		t.Errorf("expected default state_file, got %q", cfg.StateFile)
	}
	if cfg.OutputFmt != "text" {
		t.Errorf("expected default output_format 'text', got %q", cfg.OutputFmt)
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	path := writeTempConfig(t, `
state_file: custom.tfstate
region: eu-west-1
output_format: json
filters:
  - aws_instance
ignore:
  - tags
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StateFile != "custom.tfstate" {
		t.Errorf("state_file: got %q", cfg.StateFile)
	}
	if cfg.Region != "eu-west-1" {
		t.Errorf("region: got %q", cfg.Region)
	}
	if cfg.OutputFmt != "json" {
		t.Errorf("output_format: got %q", cfg.OutputFmt)
	}
	if len(cfg.Filters) != 1 || cfg.Filters[0] != "aws_instance" {
		t.Errorf("filters: got %v", cfg.Filters)
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	path := writeTempConfig(t, "output_format: xml\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for unsupported output_format")
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	path := writeTempConfig(t, ": bad: yaml: [")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected parse error for malformed YAML")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Region != "us-east-1" {
		t.Errorf("default region: got %q", cfg.Region)
	}
}
