package aws

// ResourceType represents a supported AWS resource type.
type ResourceType string

const (
	ResourceTypeS3Bucket      ResourceType = "aws_s3_bucket"
	ResourceTypeEC2Instance   ResourceType = "aws_instance"
	ResourceTypeSecurityGroup ResourceType = "aws_security_group"
	ResourceTypeIAMRole       ResourceType = "aws_iam_role"
	ResourceTypeDynamoDBTable  ResourceType = "aws_dynamodb_table"
)

// SupportedResourceTypes returns the list of AWS resource types
// that drift-check is capable of fetching and comparing.
func SupportedResourceTypes() []ResourceType {
	return []ResourceType{
		ResourceTypeS3Bucket,
		ResourceTypeEC2Instance,
		ResourceTypeSecurityGroup,
		ResourceTypeIAMRole,
		ResourceTypeDynamoDBTable,
	}
}

// IsSupported returns true if the given resource type string
// corresponds to a known, supported AWS resource type.
func IsSupported(resourceType string) bool {
	for _, rt := range SupportedResourceTypes() {
		if string(rt) == resourceType {
			return true
		}
	}
	return false
}

// LiveResource represents a resource fetched from the live AWS environment.
type LiveResource struct {
	ID         string
	Type       ResourceType
	Attributes map[string]string
}

// AttributeKeys returns the sorted list of attribute keys present
// on this live resource.
func (r *LiveResource) AttributeKeys() []string {
	keys := make([]string, 0, len(r.Attributes))
	for k := range r.Attributes {
		keys = append(keys, k)
	}
	return keys
}
