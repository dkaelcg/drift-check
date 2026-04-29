package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/drift-check/internal/aws"
	"github.com/yourorg/drift-check/internal/config"
)

var (
	scanRegion     string
	scanTypes      []string
	scanTagFilters []string
	scanMaxResults int
	scanOutput     string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan live AWS resources in a region",
	Long:  `Fetches, enriches, and filters live AWS resources then prints a summary.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		region := scanRegion
		if region == "" {
			region = cfg.AWS.Region
		}
		if region == "" {
			return fmt.Errorf("region is required (--region or config.aws.region)")
		}

		tagFilters, err := parseTagFilters(scanTagFilters)
		if err != nil {
			return fmt.Errorf("parse tag filters: %w", err)
		}

		awsCfg, err := aws.LoadConfig(context.Background(), region)
		if err != nil {
			return fmt.Errorf("aws config: %w", err)
		}

		scanner := aws.NewScanner(awsCfg)
		result, err := scanner.Scan(context.Background(), aws.ScanOptions{
			Region:     region,
			Types:      scanTypes,
			TagFilters: tagFilters,
			MaxResults: scanMaxResults,
		})
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "warning: %s\n", e)
		}

		if scanOutput == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Region : %s\n", result.Region)
		fmt.Printf("Resources found: %d\n", len(result.Resources))
		for _, r := range result.Resources {
			fmt.Printf("  [%s] %s\n", r.Type, r.ID)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVar(&scanRegion, "region", "", "AWS region to scan")
	scanCmd.Flags().StringSliceVar(&scanTypes, "type", nil, "Resource types to include (default: all supported)")
	scanCmd.Flags().StringSliceVar(&scanTagFilters, "tag", nil, "Tag filters in key=value format")
	scanCmd.Flags().IntVar(&scanMaxResults, "max-results", 0, "Limit number of returned resources (0 = unlimited)")
	scanCmd.Flags().StringVar(&scanOutput, "output", "text", "Output format: text or json")
}
