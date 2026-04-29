package aws

import (
	"strings"
	"testing"
)

func TestSuggest_EmptyDrifts(t *testing.T) {
	r := NewRemediator()
	actions := r.Suggest([]map[string]string{})
	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}

func TestSuggest_MissingResourceID(t *testing.T) {
	r := NewRemediator()
	drifts := []map[string]string{
		{"resource_type": "aws_instance", "attribute": "instance_type"},
	}
	actions := r.Suggest(drifts)
	if len(actions) != 0 {
		t.Errorf("expected 0 actions for missing resource_id, got %d", len(actions))
	}
}

func TestSuggest_WithAttribute(t *testing.T) {
	r := NewRemediator()
	drifts := []map[string]string{
		{"resource_id": "i-abc123", "resource_type": "aws_instance", "attribute": "instance_type"},
	}
	actions := r.Suggest(drifts)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if !strings.Contains(actions[0].Description, "instance_type") {
		t.Errorf("expected description to mention attribute, got: %s", actions[0].Description)
	}
	if !strings.Contains(actions[0].TFCommand, "terraform apply -target=") {
		t.Errorf("expected TFCommand to contain apply -target, got: %s", actions[0].TFCommand)
	}
}

func TestSuggest_MissingResource(t *testing.T) {
	r := NewRemediator()
	drifts := []map[string]string{
		{"resource_id": "i-missing", "resource_type": "aws_instance"},
	}
	actions := r.Suggest(drifts)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if !strings.Contains(actions[0].Description, "missing from state") {
		t.Errorf("expected missing resource description, got: %s", actions[0].Description)
	}
}

func TestSuggest_IDSanitization(t *testing.T) {
	r := NewRemediator()
	drifts := []map[string]string{
		{"resource_id": "arn:aws:s3:::my-bucket", "resource_type": "aws_s3_bucket"},
	}
	actions := r.Suggest(drifts)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if strings.ContainsAny(actions[0].TFCommand, ":-") {
		t.Errorf("TFCommand should not contain raw special chars, got: %s", actions[0].TFCommand)
	}
}

func TestSuggest_MultipleActions(t *testing.T) {
	r := NewRemediator()
	drifts := []map[string]string{
		{"resource_id": "i-aaa", "resource_type": "aws_instance", "attribute": "ami"},
		{"resource_id": "sg-bbb", "resource_type": "aws_security_group", "attribute": "description"},
	}
	actions := r.Suggest(drifts)
	if len(actions) != 2 {
		t.Errorf("expected 2 actions, got %d", len(actions))
	}
}
