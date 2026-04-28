package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the local resource snapshot cache",
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all cached resource snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("cache-dir")
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			fmt.Println("cache is already empty")
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading cache dir: %w", err)
		}
		removed := 0
		for _, e := range entries {
			if err := os.Remove(dir + "/" + e.Name()); err == nil {
				removed++
			}
		}
		fmt.Printf("removed %d cached snapshot(s)\n", removed)
		return nil
	},
}

var cacheListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cached resource snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("cache-dir")
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			fmt.Println("no cache entries found")
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading cache dir: %w", err)
		}
		for _, e := range entries {
			fmt.Println(e.Name())
		}
		return nil
	},
}

func init() {
	cacheCmd.PersistentFlags().String("cache-dir", ".drift-cache", "directory used for resource snapshot cache")
	cacheCmd.AddCommand(cacheClearCmd)
	cacheCmd.AddCommand(cacheListCmd)
	rootCmd.AddCommand(cacheCmd)
}
