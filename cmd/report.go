package cmd

import (
	"fmt"
	"os"

	"github.com/drift-check/internal/report"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Output drift results in a specified format",
	Long:  `Generate a drift report from the last detected drift results in text or JSON format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := report.NewFormatter(report.Format(outputFormat))
		if err != nil {
			return fmt.Errorf("invalid output format: %w", err)
		}

		// In a full implementation, results would be passed from the detect step.
		// Here we demonstrate the formatter wiring with an empty result set.
		results := globalDriftResults

		if err := f.Write(results, os.Stdout); err != nil {
			return fmt.Errorf("failed to write report: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(
		&outputFormat,
		"format", "f",
		"text",
		`Output format: "text" or "json"`,
	)
}
