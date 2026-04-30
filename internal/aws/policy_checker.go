package aws

import (
	"fmt"
	"strings"
)

// PolicyViolation represents a detected IAM/resource policy issue.
type PolicyViolation struct {
	ResourceID string
	ResourceType string
	Severity string
	Description string
}

// PolicyChecker evaluates resources for common policy drift or misconfigurations.
type PolicyChecker struct {
	rules []policyRule
}

type policyRule struct {
	resourceType string
	check func(resource LiveResource) *PolicyViolation
}

// NewPolicyChecker creates a PolicyChecker with built-in rules.
func NewPolicyChecker() *PolicyChecker {
	pc := &PolicyChecker{}
	pc.rules = []policyRule{
		{
			resourceType: "aws_s3_bucket",
			check: checkS3PublicAccess,
		},
		{
			resourceType: "aws_security_group",
			check: checkSecurityGroupOpenIngress,
		},
		{
			resourceType: "aws_iam_role",
			check: checkIAMRoleWildcard,
		},
	}
	return pc
}

// Check evaluates a slice of resources and returns all violations found.
func (pc *PolicyChecker) Check(resources []LiveResource) []PolicyViolation {
	var violations []PolicyViolation
	for _, res := range resources {
		for _, rule := range pc.rules {
			if rule.resourceType == res.Type {
				if v := rule.check(res); v != nil {
					violations = append(violations, *v)
				}
			}
		}
	}
	return violations
}

func checkS3PublicAccess(r LiveResource) *PolicyViolation {
	if val, ok := r.Attributes["acl"]; ok && strings.EqualFold(fmt.Sprintf("%v", val), "public-read") {
		return &PolicyViolation{
			ResourceID: r.ID,
			ResourceType: r.Type,
			Severity: "HIGH",
			Description: "S3 bucket has public-read ACL set",
		}
	}
	return nil
}

func checkSecurityGroupOpenIngress(r LiveResource) *PolicyViolation {
	if val, ok := r.Attributes["ingress_cidr"]; ok && strings.Contains(fmt.Sprintf("%v", val), "0.0.0.0/0") {
		return &PolicyViolation{
			ResourceID: r.ID,
			ResourceType: r.Type,
			Severity: "CRITICAL",
			Description: "Security group allows unrestricted ingress from 0.0.0.0/0",
		}
	}
	return nil
}

func checkIAMRoleWildcard(r LiveResource) *PolicyViolation {
	if val, ok := r.Attributes["policy"]; ok && strings.Contains(fmt.Sprintf("%v", val), "\"*\"") {
		return &PolicyViolation{
			ResourceID: r.ID,
			ResourceType: r.Type,
			Severity: "HIGH",
			Description: "IAM role policy contains wildcard (*) action or resource",
		}
	}
	return nil
}
