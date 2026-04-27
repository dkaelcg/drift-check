package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	statePath  string
	region     string
	profile    string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "drift-check",
	Short: "Detect configuration drift between live cloud resources and Terraform state",
	Long: `drift-check compares live cloud infrastructure against Terraform state files
to identify resources that have drifted from their declared configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "path to Terraform state file")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region to use (overrides profile default)")
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "AWS credentials profile")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
