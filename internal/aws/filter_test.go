package aws

import (
	"testing"
)

func sampleResources() []LiveResource {
	return []LiveResource{
		{
			ID:     "i-001",
			Type:   "aws_instance",
			Region: "us-east-1",
			Tags:   map[string]string{"env": "prod", "team": "platform"},
		},
		{
			ID:     "bucket-001",
			Type:   "aws_s3_bucket",
			Region: "us-west-2",
			Tags:   map[string]string{"env": "staging"},
		},
		{
			ID:     "sg-001",
			Type:   "aws_security_group",
			Region: "us-east-1",
			Tags:   map[string]string{},
		},
	}
}

func TestApply_NoFilter(t *testing.T) {
	res := Apply(sampleResources(), ResourceFilter{})
	if len(res.Included) != 3 {
		t.Fatalf("expected 3 included, got %d", len(res.Included))
	}
	if len(res.Excluded) != 0 {
		t.Fatalf("expected 0 excluded, got %d", len(res.Excluded))
	}
}

func TestApply_FilterByType(t *testing.T) {
	res := Apply(sampleResources(), ResourceFilter{Types: []string{"aws_instance"}})
	if len(res.Included) != 1 {
		t.Fatalf("expected 1 included, got %d", len(res.Included))
	}
	if res.Included[0].ID != "i-001" {
		t.Errorf("unexpected resource ID: %s", res.Included[0].ID)
	}
}

func TestApply_FilterByRegion(t *testing.T) {
	res := Apply(sampleResources(), ResourceFilter{Region: "us-east-1"})
	if len(res.Included) != 2 {
		t.Fatalf("expected 2 included, got %d", len(res.Included))
	}
}

func TestApply_FilterByTags(t *testing.T) {
	res := Apply(sampleResources(), ResourceFilter{
		Tags: map[string]string{"env": "prod"},
	})
	if len(res.Included) != 1 {
		t.Fatalf("expected 1 included, got %d", len(res.Included))
	}
	if res.Included[0].ID != "i-001" {
		t.Errorf("unexpected resource ID: %s", res.Included[0].ID)
	}
}

func TestApply_CombinedFilter(t *testing.T) {
	res := Apply(sampleResources(), ResourceFilter{
		Types:  []string{"aws_instance"},
		Region: "us-east-1",
		Tags:   map[string]string{"team": "platform"},
	})
	if len(res.Included) != 1 {
		t.Fatalf("expected 1 included, got %d", len(res.Included))
	}
}

func TestApply_NoMatch(t *testing.T) {
	res := Apply(sampleResources(), ResourceFilter{
		Types: []string{"aws_lambda_function"},
	})
	if len(res.Included) != 0 {
		t.Fatalf("expected 0 included, got %d", len(res.Included))
	}
	if len(res.Excluded) != 3 {
		t.Fatalf("expected 3 excluded, got %d", len(res.Excluded))
	}
}
