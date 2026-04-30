package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/drift-check/internal/aws"
	"github.com/user/drift-check/internal/report"
)

var (
	policyOutputFormat string
	policyRegion string
	policyResourceType string
)

var policyCmd = &cobra.Command{
	Use: "policy",
	Short: "Check live resources for policy violations and misconfigurations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if policyRegion == "" {
			return fmt.Errorf("--region is required")
		}

		fetcher := aws.NewFetcher(policyRegion)
		checker := aws.NewPolicyChecker()

		types := aws.SupportedResourceTypes()
		if policyResourceType != "" {
			if !aws.IsSupported(policyResourceType) {
				return fmt.Errorf("unsupported resource type: %s", policyResourceType)
			}
			types = []string{policyResourceType}
		}

		var resources []aws.LiveResource
		for _, t := range types {
			res, err := fetcher.FetchResource(t, "")
			if err != nil {
				continue
			}
			resources = append(resources, res)
		}

		violations := checker.Check(resources)
		policyReport := report.BuildPolicyReport(violations)

		switch policyOutputFormat {
		case "json":
			if err := report.WritePolicyJSON(os.Stdout, policyReport); err != nil {
				return fmt.Errorf("failed to write JSON report: %w", err)
			}
		default:
			report.WritePolicyText(os.Stdout, policyReport)
		}
		return nil
	},
}

func init() {
	policyCmd.Flags().StringVar(&policyRegion, "region", "", "AWS region to scan (required)")
	policyCmd.Flags().StringVar(&policyOutputFormat, "output", "text", "Output format: text or json")
	policyCmd.Flags().StringVar(&policyResourceType, "type", "", "Limit check to a specific resource type")
	rootCmd.AddCommand(policyCmd)
}
