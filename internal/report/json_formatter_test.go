package report_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/drift-check/internal/drift"
	"github.com/drift-check/internal/report"
)

func TestJSONFormatter_NoDrift(t *testing.T) {
	f := &report.JSONFormatter{}
	var buf bytes.Buffer
	if err := f.Write([]drift.Result{}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out["drift_detected"] != false {
		t.Errorf("expected drift_detected=false, got: %v", out["drift_detected"])
	}
}

func TestJSONFormatter_WithDrift(t *testing.T) {
	f := &report.JSONFormatter{}
	var buf bytes.Buffer
	if err := f.Write(driftResults(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out["drift_detected"] != true {
		t.Errorf("expected drift_detected=true, got: %v", out["drift_detected"])
	}

	resources, ok := out["resources"].([]interface{})
	if !ok || len(resources) == 0 {
		t.Fatal("expected at least one resource in output")
	}

	res := resources[0].(map[string]interface{})
	if res["id"] != "i-abc123" {
		t.Errorf("expected resource id 'i-abc123', got: %v", res["id"])
	}
	if res["has_drift"] != true {
		t.Errorf("expected has_drift=true, got: %v", res["has_drift"])
	}
}

func TestJSONFormatter_DifferencesIncluded(t *testing.T) {
	f := &report.JSONFormatter{}
	var buf bytes.Buffer
	if err := f.Write(driftResults(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	json.Unmarshal(buf.Bytes(), &out)
	resources := out["resources"].([]interface{})
	res := resources[0].(map[string]interface{})
	diffs, ok := res["differences"].([]interface{})
	if !ok || len(diffs) == 0 {
		t.Fatal("expected differences in output")
	}
}
