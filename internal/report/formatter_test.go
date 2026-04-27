package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/drift-check/internal/drift"
	"github.com/drift-check/internal/report"
)

func driftResults() []drift.Result {
	return []drift.Result{
		{
			ResourceID:   "i-abc123",
			ResourceType: "aws_instance",
			Differences: []drift.Difference{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t2.small"},
			},
		},
	}
}

func TestNewFormatter_Text(t *testing.T) {
	f, err := report.NewFormatter(report.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestNewFormatter_JSON(t *testing.T) {
	f, err := report.NewFormatter(report.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestNewFormatter_Unsupported(t *testing.T) {
	_, err := report.NewFormatter("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestTextFormatter_NoDrift(t *testing.T) {
	f := &report.TextFormatter{}
	var buf bytes.Buffer
	if err := f.Write(nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestTextFormatter_WithDrift(t *testing.T) {
	f := &report.TextFormatter{}
	var buf bytes.Buffer
	if err := f.Write(driftResults(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "i-abc123") {
		t.Errorf("expected resource ID in output, got: %s", out)
	}
	if !strings.Contains(out, "instance_type") {
		t.Errorf("expected attribute name in output, got: %s", out)
	}
}
