package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/drift-check/internal/config"
)

var configPath string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display the resolved drift-check configuration",
	Long: `Load and validate the drift-check YAML configuration file,
then print the resolved settings that will be used during a run.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Configuration resolved from: %s\n", configPath)
		fmt.Fprintf(cmd.OutOrStdout(), "  state_file    : %s\n", cfg.StateFile)
		fmt.Fprintf(cmd.OutOrStdout(), "  region        : %s\n", cfg.Region)
		fmt.Fprintf(cmd.OutOrStdout(), "  output_format : %s\n", cfg.OutputFmt)

		if len(cfg.Filters) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  filters       : %v\n", cfg.Filters)
		}
		if len(cfg.Ignore) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  ignore        : %v\n", cfg.Ignore)
		}

		return nil
	},
}

func init() {
	configCmd.Flags().StringVarP(
		&configPath, "file", "f", "drift.yaml",
		"path to the drift-check YAML configuration file",
	)
	rootCmd.AddCommand(configCmd)
}
