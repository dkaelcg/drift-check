package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/drift-check/internal/aws"
)

var (
	remediateOutput string
)

var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Suggest remediation commands for detected drift",
	Long:  `Reads a JSON drift report and outputs terraform commands to resolve each drift.`,
	RunE:  runRemediate,
}

func runRemediate(cmd *cobra.Command, args []string) error {
	inputFile, _ := cmd.Flags().GetString("input")
	if inputFile == "" {
		return fmt.Errorf("--input flag is required")
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	var drifts []map[string]string
	if err := json.Unmarshal(data, &drifts); err != nil {
		return fmt.Errorf("failed to parse drift input: %w", err)
	}

	remediator := aws.NewRemediator()
	actions := remediator.Suggest(drifts)

	if len(actions) == 0 {
		fmt.Println("No remediation actions needed.")
		return nil
	}

	if remediateOutput == "json" {
		out, err := json.MarshalIndent(actions, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal actions: %w", err)
		}
		fmt.Println(string(out))
		return nil
	}

	for _, a := range actions {
		fmt.Printf("[%s] %s\n  => %s\n\n", a.ResourceType, a.Description, a.TFCommand)
	}
	return nil
}

func init() {
	remediateCmd.Flags().String("input", "", "Path to JSON drift report file (required)")
	remediateCmd.Flags().StringVarP(&remediateOutput, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(remediateCmd)
}
