//go:build integration
// +build integration

package aws

import (
	"context"
	"os"
	"testing"
)

// TestScanner_Integration_LiveScan runs a real AWS scan.
// Requires AWS credentials and the 'integration' build tag.
func TestScanner_Integration_LiveScan(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		t.Skip("AWS_REGION not set, skipping integration test")
	}

	cfg, err := LoadConfig(context.Background(), region)
	if err != nil {
		t.Fatalf("load aws config: %v", err)
	}

	scanner := NewScanner(cfg)
	result, err := scanner.Scan(context.Background(), ScanOptions{
		Region:     region,
		MaxResults: 10,
	})
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	t.Logf("region=%s resources=%d warnings=%d",
		result.Region, len(result.Resources), len(result.Errors))

	for _, e := range result.Errors {
		t.Logf("warning: %s", e)
	}

	for _, r := range result.Resources {
		if r.ID == "" {
			t.Errorf("resource of type %s has empty ID", r.Type)
		}
		if r.Type == "" {
			t.Error("resource has empty type")
		}
	}
}
