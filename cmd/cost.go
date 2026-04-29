package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/drift-check/internal/aws"
)

var costOutputFormat string

var costCmd = &cobra.Command{
	Use:   "cost",
	Short: "Estimate monthly costs for discovered live resources",
	Long: `Scans live AWS resources and provides rough monthly cost estimates
based on resource type. Useful for identifying potentially expensive
resources that may have drifted from their expected configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		region, _ := cmd.Flags().GetString("region")
		resType, _ := cmd.Flags().GetString("type")

		if region == "" {
			return fmt.Errorf("--region flag is required")
		}

		scanner := aws.NewScanner(aws.ScannerConfig{
			Region:     region,
			MaxResults: 100,
		})

		resources, err := scanner.Scan(resType)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		estimator := aws.NewCostEstimator()
		estimates := estimator.EstimateAll(resources)
		total := aws.TotalMonthlyCost(estimates)

		if costOutputFormat == "json" {
			output := map[string]interface{}{
				"estimates":           estimates,
				"total_monthly_cost":  total,
				"currency":            "USD",
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(output)
		}

		fmt.Printf("Cost estimates for region: %s\n\n", region)
		for _, est := range estimates {
			line := fmt.Sprintf("  %-40s %-30s $%8.2f/mo", est.ResourceID, est.ResourceType, est.MonthlyCost)
			if est.Note != "" {
				line += fmt.Sprintf(" (%s)", est.Note)
			}
			fmt.Println(line)
		}
		fmt.Printf("\nTotal estimated monthly cost: $%.2f USD\n", total)
		return nil
	},
}

func init() {
	costCmd.Flags().String("region", "", "AWS region to scan (required)")
	costCmd.Flags().String("type", "", "Filter by resource type (optional)")
	costCmd.Flags().StringVar(&costOutputFormat, "output", "text", "Output format: text or json")
	rootCmd.AddCommand(costCmd)
}
