package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/drift-check/internal/drift"
)

// Format defines the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes a drift report to an io.Writer.
type Formatter interface {
	Write(results []drift.Result, w io.Writer) error
}

// NewFormatter returns a Formatter for the given format.
func NewFormatter(f Format) (Formatter, error) {
	switch f {
	case FormatText:
		return &TextFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %q", f)
	}
}

// TextFormatter writes a human-readable drift report.
type TextFormatter struct{}

func (t *TextFormatter) Write(results []drift.Result, w io.Writer) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "✓ No drift detected.")
		return err
	}

	for _, r := range results {
		if !r.HasDrift() {
			continue
		}
		fmt.Fprintf(w, "[DRIFT] Resource: %s (%s)\n", r.ResourceID, r.ResourceType)
		for _, d := range r.Differences {
			fmt.Fprintf(w, "  - %s: state=%q live=%q\n",
				d.Attribute,
				strings.TrimSpace(d.StateValue),
				strings.TrimSpace(d.LiveValue),
			)
		}
	}
	return nil
}
