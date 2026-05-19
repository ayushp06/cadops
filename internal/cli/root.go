package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// NewRootCmd creates the top-level cadops command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cadops",
		Short:        "CAD-aware Git workflow helpers",
		SilenceUsage: true,
		Version:      versionString(),
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newDiffCmd())
	cmd.AddCommand(newFilesCmd())
	cmd.AddCommand(newScanCmd())
	cmd.AddCommand(newMetadataCmd())
	cmd.AddCommand(newPreviewCmd())
	cmd.AddCommand(newDoctorCmd())
	cmd.AddCommand(newWatchCmd())
	cmd.AddCommand(newSnapshotCmd())
	cmd.AddCommand(newCommitCmd())
	cmd.AddCommand(newLockCmd())
	cmd.AddCommand(newUnlockCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newPushCmd())
	cmd.AddCommand(newPullCmd())
	cmd.AddCommand(newHistoryCmd())
	cmd.AddCommand(newVersionCmd())

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CadOps version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), versionString())
		},
	}
}

func versionString() string {
	return fmt.Sprintf("%s (commit %s, built %s)", version, commit, date)
}
