package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/example/drift-check/internal/snapshot"
	"github.com/spf13/cobra"
)

var snapshotDir string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage drift detection snapshots",
}

var snapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := snapshot.NewManager(snapshotDir)
		if err != nil {
			return err
		}
		ids, err := m.List()
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			fmt.Println("No snapshots found.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATE FILE\tCREATED AT")
		for _, id := range ids {
			s, err := m.Load(id)
			if err != nil {
				fmt.Fprintf(w, "%s\t(unreadable)\t-\n", id)
				continue
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", s.ID, s.StateFile, s.CreatedAt.Format(time.RFC3339))
		}
		return w.Flush()
	},
}

var snapshotDiffCmd = &cobra.Command{
	Use:   "diff <prev-id> <curr-id>",
	Short: "Compare two snapshots and show drift changes",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := snapshot.NewManager(snapshotDir)
		if err != nil {
			return err
		}
		prev, err := m.Load(args[0])
		if err != nil {
			return fmt.Errorf("load prev snapshot: %w", err)
		}
		curr, err := m.Load(args[1])
		if err != nil {
			return fmt.Errorf("load curr snapshot: %w", err)
		}
		result, err := snapshot.Compare(prev, curr)
		if err != nil {
			return err
		}
		fmt.Printf("Added:   %v\n", result.Added)
		fmt.Printf("Removed: %v\n", result.Removed)
		for _, rd := range result.Changed {
			fmt.Printf("Changed: %s\n", rd.ResourceID)
			for _, fd := range rd.Changes {
				fmt.Printf("  %s: %q -> %q\n", fd.Field, fd.Previous, fd.Current)
			}
		}
		return nil
	},
}

func init() {
	snapshotCmd.PersistentFlags().StringVar(&snapshotDir, "dir", ".drift-snapshots", "Directory to store snapshots")
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotDiffCmd)
	rootCmd.AddCommand(snapshotCmd)
}
