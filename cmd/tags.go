package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/drift-check/internal/aws"
)

var (
	tagFiltersFlag []string
	tagOutputJSON  bool
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Fetch AWS resources by tag filters",
	Long:  `Query AWS Resource Groups Tagging API to list resources matching one or more tag filters (KEY=VALUE).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filters, err := parseTagFilters(tagFiltersFlag)
		if err != nil {
			return fmt.Errorf("invalid tag filter: %w", err)
		}

		cfg, err := loadAWSConfig(cmd.Context())
		if err != nil {
			return fmt.Errorf("aws config: %w", err)
		}

		client := aws.NewTaggingClientFromConfig(cfg)
		tagger := aws.NewTagger(client)

		resources, err := tagger.FetchByTags(context.Background(), filters)
		if err != nil {
			return fmt.Errorf("fetch by tags: %w", err)
		}

		if tagOutputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(resources)
		}

		if len(resources) == 0 {
			fmt.Println("No resources found matching the given tag filters.")
			return nil
		}

		for _, r := range resources {
			fmt.Printf("ARN: %s\n", r.ARN)
			for k, v := range r.Tags {
				fmt.Printf("  %s = %s\n", k, v)
			}
		}
		return nil
	},
}

func parseTagFilters(raw []string) ([]aws.TagFilter, error) {
	var filters []aws.TagFilter
	for _, entry := range raw {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("expected KEY=VALUE, got %q", entry)
		}
		filters = append(filters, aws.TagFilter{
			Key:    parts[0],
			Values: []string{parts[1]},
		})
	}
	return filters, nil
}

func init() {
	tafsCmd := tagsCmd
	tafsCmd.Flags().StringArrayVarP(&tagFiltersFlag, "filter", "f", nil, "Tag filter in KEY=VALUE format (repeatable)")
	tafsCmd.Flags().BoolVar(&tagOutputJSON, "json", false, "Output results as JSON")
	rootCmd.AddCommand(tafsCmd)
}
