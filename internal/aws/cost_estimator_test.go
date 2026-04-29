package aws

import (
	"testing"
)

func TestEstimate_KnownType(t *testing.T) {
	e := NewCostEstimator()
	r := LiveResource{ID: "i-abc123", Type: "aws_instance", Region: "us-east-1"}

	est, err := e.Estimate(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if est.MonthlyCost != 72.00 {
		t.Errorf("expected 72.00, got %.2f", est.MonthlyCost)
	}
	if est.Currency != "USD" {
		t.Errorf("expected USD, got %s", est.Currency)
	}
	if est.Note != "" {
		t.Errorf("expected empty note for known type, got %q", est.Note)
	}
}

func TestEstimate_UnknownType(t *testing.T) {
	e := NewCostEstimator()
	r := LiveResource{ID: "res-xyz", Type: "aws_unknown_resource", Region: "eu-west-1"}

	est, err := e.Estimate(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if est.MonthlyCost != 0.0 {
		t.Errorf("expected 0.0 for unknown type, got %.2f", est.MonthlyCost)
	}
	if est.Note == "" {
		t.Error("expected a note for unknown resource type")
	}
}

func TestEstimate_EmptyID(t *testing.T) {
	e := NewCostEstimator()
	r := LiveResource{ID: "", Type: "aws_instance", Region: "us-east-1"}

	_, err := e.Estimate(r)
	if err == nil {
		t.Fatal("expected error for empty resource ID")
	}
}

func TestEstimateAll_MultipleResources(t *testing.T) {
	e := NewCostEstimator()
	resources := []LiveResource{
		{ID: "i-1", Type: "aws_instance", Region: "us-east-1"},
		{ID: "bucket-1", Type: "aws_s3_bucket", Region: "us-east-1"},
		{ID: "", Type: "aws_eip", Region: "us-east-1"}, // skipped
	}

	estimates := e.EstimateAll(resources)
	if len(estimates) != 2 {
		t.Errorf("expected 2 estimates (empty ID skipped), got %d", len(estimates))
	}
}

func TestTotalMonthlyCost(t *testing.T) {
	estimates := []CostEstimate{
		{MonthlyCost: 72.00},
		{MonthlyCost: 2.30},
		{MonthlyCost: 0.0},
	}

	total := TotalMonthlyCost(estimates)
	expected := 74.30
	if total != expected {
		t.Errorf("expected %.2f, got %.2f", expected, total)
	}
}

func TestTotalMonthlyCost_Empty(t *testing.T) {
	total := TotalMonthlyCost([]CostEstimate{})
	if total != 0.0 {
		t.Errorf("expected 0.0 for empty estimates, got %.2f", total)
	}
}
