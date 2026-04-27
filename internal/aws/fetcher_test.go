package aws

import (
	"testing"
)

func TestLiveResource_Fields(t *testing.T) {
	r := &LiveResource{
		Type: "aws_instance",
		ID:   "i-0abc123",
		Attributes: ResourceAttributes{
			"instance_type": "t3.micro",
			"instance_state": "running",
			"ami":            "ami-0deadbeef",
		},
	}

	if r.Type != "aws_instance" {
		t.Errorf("expected type aws_instance, got %s", r.Type)
	}
	if r.ID != "i-0abc123" {
		t.Errorf("expected ID i-0abc123, got %s", r.ID)
	}
	if r.Attributes["instance_type"] != "t3.micro" {
		t.Errorf("expected instance_type t3.micro, got %s", r.Attributes["instance_type"])
	}
}

func TestFetchResource_UnsupportedType(t *testing.T) {
	// Fetcher with zero-value config to test dispatch without real AWS calls.
	f := &Fetcher{}
	_, err := f.FetchResource(nil, "aws_lambda_function", "my-func") //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for unsupported resource type, got nil")
	}
	expected := "unsupported resource type: aws_lambda_function"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestResourceAttributes_Empty(t *testing.T) {
	attrs := ResourceAttributes{}
	if len(attrs) != 0 {
		t.Errorf("expected empty attributes map")
	}
	attrs["key"] = "value"
	if attrs["key"] != "value" {
		t.Errorf("expected value 'value', got %s", attrs["key"])
	}
}
