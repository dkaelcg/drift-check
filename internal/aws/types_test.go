package aws

import (
	"testing"
)

func TestSupportedResourceTypes_NotEmpty(t *testing.T) {
	types := SupportedResourceTypes()
	if len(types) == 0 {
		t.Fatal("expected at least one supported resource type, got none")
	}
}

func TestSupportedResourceTypes_ContainsKnown(t *testing.T) {
	expected := []ResourceType{
		ResourceTypeS3Bucket,
		ResourceTypeEC2Instance,
		ResourceTypeSecurityGroup,
		ResourceTypeIAMRole,
		ResourceTypeDynamoDBTable,
	}
	types := SupportedResourceTypes()
	typeSet := make(map[ResourceType]bool, len(types))
	for _, rt := range types {
		typeSet[rt] = true
	}
	for _, e := range expected {
		if !typeSet[e] {
			t.Errorf("expected resource type %q to be supported", e)
		}
	}
}

func TestIsSupported_KnownType(t *testing.T) {
	if !IsSupported("aws_s3_bucket") {
		t.Error("expected aws_s3_bucket to be supported")
	}
}

func TestIsSupported_UnknownType(t *testing.T) {
	if IsSupported("aws_unknown_resource") {
		t.Error("expected aws_unknown_resource to be unsupported")
	}
}

func TestIsSupported_EmptyString(t *testing.T) {
	if IsSupported("") {
		t.Error("expected empty string to be unsupported")
	}
}

func TestLiveResource_AttributeKeys_Empty(t *testing.T) {
	r := &LiveResource{
		ID:         "test-id",
		Type:       ResourceTypeS3Bucket,
		Attributes: map[string]string{},
	}
	keys := r.AttributeKeys()
	if len(keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(keys))
	}
}

func TestLiveResource_AttributeKeys_NonEmpty(t *testing.T) {
	r := &LiveResource{
		ID:   "i-12345",
		Type: ResourceTypeEC2Instance,
		Attributes: map[string]string{
			"instance_type": "t3.micro",
			"ami":           "ami-abc123",
		},
	}
	keys := r.AttributeKeys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
