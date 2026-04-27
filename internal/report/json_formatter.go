package report

import (
	"encoding/json"
	"io"

	"github.com/drift-check/internal/drift"
)

// JSONFormatter writes a machine-readable JSON drift report.
type JSONFormatter struct{}

type jsonReport struct {
	DriftDetected bool          `json:"drift_detected"`
	Resources     []jsonResource `json:"resources"`
}

type jsonResource struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	HasDrift    bool            `json:"has_drift"`
	Differences []jsonDiff      `json:"differences,omitempty"`
}

type jsonDiff struct {
	Attribute  string `json:"attribute"`
	StateValue string `json:"state_value"`
	LiveValue  string `json:"live_value"`
}

func (j *JSONFormatter) Write(results []drift.Result, w io.Writer) error {
	report := jsonReport{
		Resources: make([]jsonResource, 0, len(results)),
	}

	for _, r := range results {
		jr := jsonResource{
			ID:       r.ResourceID,
			Type:     r.ResourceType,
			HasDrift: r.HasDrift(),
		}
		for _, d := range r.Differences {
			jr.Differences = append(jr.Differences, jsonDiff{
				Attribute:  d.Attribute,
				StateValue: d.StateValue,
				LiveValue:  d.LiveValue,
			})
		}
		report.Resources = append(report.Resources, jr)
		if r.HasDrift() {
			report.DriftDetected = true
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
