package aws

import (
	"testing"
)

func baseResource() LiveResource {
	return LiveResource{
		ID:         "i-0abc123",
		Type:       "aws_instance",
		Attributes: map[string]string{"instance_type": "t3.micro"},
	}
}

func TestEnrich_Basic(t *testing.T) {
	e := NewEnricher("us-east-1", "123456789012")
	r := baseResource()
	tags := map[string]string{"Env": "prod", "Team": "platform"}

	enriched, err := e.Enrich(r, tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enriched.Region != "us-east-1" {
		t.Errorf("expected region us-east-1, got %s", enriched.Region)
	}
	if enriched.AccountID != "123456789012" {
		t.Errorf("expected account 123456789012, got %s", enriched.AccountID)
	}
	if len(enriched.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(enriched.Tags))
	}
}

func TestEnrich_EmptyID(t *testing.T) {
	e := NewEnricher("eu-west-1", "000000000000")
	r := LiveResource{ID: "", Type: "aws_instance"}

	_, err := e.Enrich(r, nil)
	if err == nil {
		t.Fatal("expected error for empty resource ID, got nil")
	}
}

func TestEnrich_NoTags(t *testing.T) {
	e := NewEnricher("us-west-2", "111111111111")
	r := baseResource()

	enriched, err := e.Enrich(r, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(enriched.Tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(enriched.Tags))
	}
}

func TestEnrich_MetadataSource(t *testing.T) {
	e := NewEnricher("ap-southeast-1", "222222222222")
	r := baseResource()

	enriched, err := e.Enrich(r, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enriched.Metadata["source"] != "aws-live" {
		t.Errorf("expected metadata source 'aws-live', got %q", enriched.Metadata["source"])
	}
}

func TestTagMap_RoundTrip(t *testing.T) {
	e := NewEnricher("us-east-2", "333333333333")
	r := baseResource()
	tags := map[string]string{"Name": "web-server", "CostCenter": "42"}

	enriched, err := e.Enrich(r, tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := enriched.TagMap()
	for k, v := range tags {
		if result[k] != v {
			t.Errorf("tag mismatch for key %q: want %q got %q", k, v, result[k])
		}
	}
}
