package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the top-level cadops command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cadops",
		Short:        "CAD-aware Git workflow helpers",
		SilenceUsage: true,
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newDoctorCmd())
	cmd.AddCommand(newWatchCmd())
	cmd.AddCommand(newSnapshotCmd())
	cmd.AddCommand(newLockCmd())
	cmd.AddCommand(newUnlockCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newPushCmd())
	cmd.AddCommand(newPullCmd())
	cmd.AddCommand(newHistoryCmd())

	return cmd
}
