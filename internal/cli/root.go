package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the top-level cadops command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cadops",
		Short: "CAD-aware Git workflow helpers",
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newDoctorCmd())

	return cmd
}
