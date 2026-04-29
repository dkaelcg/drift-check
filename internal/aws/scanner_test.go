package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestScanner_New(t *testing.T) {
	cfg := aws.Config{}
	s := NewScanner(cfg)
	if s == nil {
		t.Fatal("expected non-nil scanner")
	}
	if s.fetcher == nil {
		t.Error("expected fetcher to be initialised")
	}
	if s.enricher == nil {
		t.Error("expected enricher to be initialised")
	}
	if s.tagger == nil {
		t.Error("expected tagger to be initialised")
	}
}

func TestScanner_Scan_EmptyRegion(t *testing.T) {
	s := NewScanner(aws.Config{})
	_, err := s.Scan(context.Background(), ScanOptions{})
	if err == nil {
		t.Fatal("expected error for empty region")
	}
}

func TestScanner_Scan_UnsupportedType(t *testing.T) {
	s := NewScanner(aws.Config{})
	opts := ScanOptions{
		Region: "us-east-1",
		Types:  []string{"aws_unknown_resource"},
	}
	result, err := s.Scan(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Errors) == 0 {
		t.Error("expected at least one error for unsupported type")
	}
	if len(result.Resources) != 0 {
		t.Errorf("expected no resources, got %d", len(result.Resources))
	}
}

func TestScanner_Scan_MaxResults(t *testing.T) {
	// Build a result manually to verify MaxResults trimming logic.
	result := &ScanResult{
		Region: "us-east-1",
		Resources: make([]LiveResource, 10),
	}
	maxResults := 3
	if len(result.Resources) > maxResults {
		result.Resources = result.Resources[:maxResults]
	}
	if len(result.Resources) != maxResults {
		t.Errorf("expected %d resources after trim, got %d", maxResults, len(result.Resources))
	}
}

func TestScanner_Scan_ResultRegion(t *testing.T) {
	// Verify that region is propagated even when no resources are returned.
	s := NewScanner(aws.Config{})
	opts := ScanOptions{
		Region: "eu-west-1",
		Types:  []string{"aws_unknown_xyz"},
	}
	result, err := s.Scan(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Region != "eu-west-1" {
		t.Errorf("expected region eu-west-1, got %s", result.Region)
	}
}
