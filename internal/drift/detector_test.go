package drift

import (
	"testing"
)

func TestDetect_NoDrift(t *testing.T) {
	state := []StateResource{
		{ID: "i-123", Type: "aws_instance", Attributes: map[string]string{"instance_type": "t3.micro", "region": "us-east-1"}},
	}
	live := map[string]LiveResource{
		"i-123": {ID: "i-123", Type: "aws_instance", Attributes: map[string]string{"instance_type": "t3.micro", "region": "us-east-1"}},
	}

	results, err := Detect(state, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].HasDrift {
		t.Errorf("expected no drift, but drift was detected")
	}
	if len(results[0].Differences) != 0 {
		t.Errorf("expected 0 differences, got %d", len(results[0].Differences))
	}
}

func TestDetect_AttributeMismatch(t *testing.T) {
	state := []StateResource{
		{ID: "i-456", Type: "aws_instance", Attributes: map[string]string{"instance_type": "t3.micro"}},
	}
	live := map[string]LiveResource{
		"i-456": {ID: "i-456", Type: "aws_instance", Attributes: map[string]string{"instance_type": "t3.large"}},
	}

	results, err := Detect(state, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].HasDrift {
		t.Errorf("expected drift to be detected")
	}
	if len(results[0].Differences) != 1 {
		t.Fatalf("expected 1 difference, got %d", len(results[0].Differences))
	}
	d := results[0].Differences[0]
	if d.Attribute != "instance_type" || d.StateValue != "t3.micro" || d.LiveValue != "t3.large" {
		t.Errorf("unexpected difference: %+v", d)
	}
}

func TestDetect_ResourceMissing(t *testing.T) {
	state := []StateResource{
		{ID: "i-789", Type: "aws_instance", Attributes: map[string]string{}},
	}
	live := map[string]LiveResource{}

	results, err := Detect(state, live)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].HasDrift {
		t.Errorf("expected drift for missing resource")
	}
	if results[0].Differences[0].Attribute != "<resource>" {
		t.Errorf("expected '<resource>' attribute marker for missing resource")
	}
}

func TestDetect_NilStateResources(t *testing.T) {
	_, err := Detect(nil, map[string]LiveResource{})
	if err == nil {
		t.Error("expected error for nil stateResources, got nil")
	}
}

func TestDetect_EmptyInputs(t *testing.T) {
	results, err := Detect([]StateResource{}, map[string]LiveResource{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
