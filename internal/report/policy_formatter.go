package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/user/drift-check/internal/aws"
)

// PolicyReport holds the output of a policy check run.
type PolicyReport struct {
	Total int `json:"total"`
	Critical int `json:"critical"`
	High int `json:"high"`
	Violations []aws.PolicyViolation `json:"violations"`
}

// BuildPolicyReport aggregates violations into a PolicyReport.
func BuildPolicyReport(violations []aws.PolicyViolation) PolicyReport {
	report := PolicyReport{Violations: violations, Total: len(violations)}
	for _, v := range violations {
		switch v.Severity {
		case "CRITICAL":
			report.Critical++
		case "HIGH":
			report.High++
		}
	}
	return report
}

// WritePolicyText writes a human-readable policy report to w.
func WritePolicyText(w io.Writer, report PolicyReport) {
	if report.Total == 0 {
		fmt.Fprintln(w, "No policy violations detected.")
		return
	}
	fmt.Fprintf(w, "Policy Violations: %d (CRITICAL: %d, HIGH: %d)\n\n",
		report.Total, report.Critical, report.High)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE ID\tTYPE\tSEVERITY\tDESCRIPTION")
	for _, v := range report.Violations {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", v.ResourceID, v.ResourceType, v.Severity, v.Description)
	}
	tw.Flush()
}

// WritePolicyJSON writes a JSON-encoded policy report to w.
func WritePolicyJSON(w io.Writer, report PolicyReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
