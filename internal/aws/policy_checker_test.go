package aws

import (
	"testing"
)

func baseCheckerResource(id, resType string, attrs map[string]interface{}) LiveResource {
	return LiveResource{
		ID: id,
		Type: resType,
		Region: "us-east-1",
		Attributes: attrs,
	}
}

func TestPolicyChecker_NoViolations(t *testing.T) {
	pc := NewPolicyChecker()
	resources := []LiveResource{
		baseCheckerResource("bucket-1", "aws_s3_bucket", map[string]interface{}{"acl": "private"}),
	}
	violations := pc.Check(resources)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestPolicyChecker_S3PublicRead(t *testing.T) {
	pc := NewPolicyChecker()
	resources := []LiveResource{
		baseCheckerResource("bucket-2", "aws_s3_bucket", map[string]interface{}{"acl": "public-read"}),
	}
	violations := pc.Check(resources)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Severity != "HIGH" {
		t.Errorf("expected HIGH severity, got %s", violations[0].Severity)
	}
}

func TestPolicyChecker_SecurityGroupOpenIngress(t *testing.T) {
	pc := NewPolicyChecker()
	resources := []LiveResource{
		baseCheckerResource("sg-1", "aws_security_group", map[string]interface{}{"ingress_cidr": "0.0.0.0/0"}),
	}
	violations := pc.Check(resources)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Severity != "CRITICAL" {
		t.Errorf("expected CRITICAL severity, got %s", violations[0].Severity)
	}
}

func TestPolicyChecker_IAMRoleWildcard(t *testing.T) {
	pc := NewPolicyChecker()
	resources := []LiveResource{
		baseCheckerResource("role-1", "aws_iam_role", map[string]interface{}{"policy": `{"Action":"*"}`}),
	}
	violations := pc.Check(resources)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestPolicyChecker_MultipleViolations(t *testing.T) {
	pc := NewPolicyChecker()
	resources := []LiveResource{
		baseCheckerResource("bucket-3", "aws_s3_bucket", map[string]interface{}{"acl": "public-read"}),
		baseCheckerResource("sg-2", "aws_security_group", map[string]interface{}{"ingress_cidr": "0.0.0.0/0"}),
	}
	violations := pc.Check(resources)
	if len(violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(violations))
	}
}

func TestPolicyChecker_UnsupportedType(t *testing.T) {
	pc := NewPolicyChecker()
	resources := []LiveResource{
		baseCheckerResource("rds-1", "aws_db_instance", map[string]interface{}{"acl": "public-read"}),
	}
	violations := pc.Check(resources)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations for unsupported type, got %d", len(violations))
	}
}
