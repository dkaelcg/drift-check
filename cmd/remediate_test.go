package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempDriftReport(t *testing.T, drifts []map[string]string) string {
	t.Helper()
	data, err := json.Marshal(drifts)
	if err != nil {
		t.Fatalf("failed to marshal drifts: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "drift-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRunRemediate_NoInput(t *testing.T) {
	cmd := remediateCmd
	cmd.ResetFlags()
	init()
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("expected error when --input is missing")
	}
}

func TestRunRemediate_InvalidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	remediateCmd.Flags().Set("input", path)
	err := remediateCmd.RunE(remediateCmd, []string{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunRemediate_MalformedJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	f.WriteString("not-json")
	f.Close()
	remediateCmd.Flags().Set("input", f.Name())
	err := remediateCmd.RunE(remediateCmd, []string{})
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestRunRemediate_ValidInput(t *testing.T) {
	drifts := []map[string]string{
		{"resource_id": "i-abc", "resource_type": "aws_instance", "attribute": "ami"},
	}
	path := writeTempDriftReport(t, drifts)
	remediateCmd.Flags().Set("input", path)
	remediateCmd.Flags().Set("output", "text")
	err := remediateCmd.RunE(remediateCmd, []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunRemediate_JSONOutput(t *testing.T) {
	drifts := []map[string]string{
		{"resource_id": "i-xyz", "resource_type": "aws_instance", "attribute": "instance_type"},
	}
	path := writeTempDriftReport(t, drifts)
	remediateCmd.Flags().Set("input", path)
	remediateCmd.Flags().Set("output", "json")
	err := remediateCmd.RunE(remediateCmd, []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
