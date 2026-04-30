package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/drift-check/internal/aws"
)

var sampleViolations = []aws.PolicyViolation{
	{ResourceID: "bucket-1", ResourceType: "aws_s3_bucket", Severity: "HIGH", Description: "public-read ACL"},
	{ResourceID: "sg-1", ResourceType: "aws_security_group", Severity: "CRITICAL", Description: "open ingress"},
}

func TestBuildPolicyReport_Counts(t *testing.T) {
	report := BuildPolicyReport(sampleViolations)
	if report.Total != 2 {
		t.Errorf("expected Total=2, got %d", report.Total)
	}
	if report.Critical != 1 {
		t.Errorf("expected Critical=1, got %d", report.Critical)
	}
	if report.High != 1 {
		t.Errorf("expected High=1, got %d", report.High)
	}
}

func TestBuildPolicyReport_Empty(t *testing.T) {
	report := BuildPolicyReport(nil)
	if report.Total != 0 {
		t.Errorf("expected Total=0, got %d", report.Total)
	}
}

func TestWritePolicyText_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	WritePolicyText(&buf, BuildPolicyReport(nil))
	if !strings.Contains(buf.String(), "No policy violations") {
		t.Errorf("expected no-violations message, got: %s", buf.String())
	}
}

func TestWritePolicyText_WithViolations(t *testing.T) {
	var buf bytes.Buffer
	WritePolicyText(&buf, BuildPolicyReport(sampleViolations))
	output := buf.String()
	if !strings.Contains(output, "bucket-1") {
		t.Errorf("expected resource ID in output")
	}
	if !strings.Contains(output, "CRITICAL") {
		t.Errorf("expected CRITICAL severity in output")
	}
}

func TestWritePolicyJSON_Valid(t *testing.T) {
	var buf bytes.Buffer
	report := BuildPolicyReport(sampleViolations)
	if err := WritePolicyJSON(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded PolicyReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if decoded.Total != 2 {
		t.Errorf("expected Total=2 in JSON, got %d", decoded.Total)
	}
}
