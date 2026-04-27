package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/drift-check/internal/state"
)

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "terraform.tfstate")
	if err := os.WriteFile(tmp, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp state: %v", err)
	}
	return tmp
}

const validState = `{
  "version": 4,
  "terraform_version": "1.5.0",
  "resources": [
    {
      "mode": "managed",
      "type": "aws_instance",
      "name": "web",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [{"schema_version": 1, "attributes": {"id": "i-abc123", "instance_type": "t3.micro"}}]
    },
    {
      "mode": "data",
      "type": "aws_ami",
      "name": "ubuntu",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [{"schema_version": 0, "attributes": {"id": "ami-xyz"}}]
    }
  ]
}`

func TestParseStateFile_Valid(t *testing.T) {
	path := writeTempState(t, validState)
	s, err := state.ParseStateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Version != 4 {
		t.Errorf("expected version 4, got %d", s.Version)
	}
	if len(s.Resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(s.Resources))
	}
}

func TestParseStateFile_ManagedOnly(t *testing.T) {
	path := writeTempState(t, validState)
	s, _ := state.ParseStateFile(path)
	managed := s.ManagedResources()
	if len(managed) != 1 {
		t.Errorf("expected 1 managed resource, got %d", len(managed))
	}
	if managed[0].Type != "aws_instance" {
		t.Errorf("expected aws_instance, got %s", managed[0].Type)
	}
}

func TestParseStateFile_UnsupportedVersion(t *testing.T) {
	path := writeTempState(t, `{"version": 3, "terraform_version": "0.12.0", "resources": []}`)
	_, err := state.ParseStateFile(path)
	if err == nil {
		t.Fatal("expected error for unsupported version, got nil")
	}
}

func TestParseStateFile_NotFound(t *testing.T) {
	_, err := state.ParseStateFile("/nonexistent/terraform.tfstate")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
