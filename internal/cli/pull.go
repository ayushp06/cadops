package cli

import (
	"fmt"
	"os"

	"github.com/cadops/cadops/internal/collab"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/spf13/cobra"
)

func newPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Run CAD-aware checks before git pull",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runPull(dir)
		},
	}
}

func runPull(dir string) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	entries, err := gitx.StatusPorcelain(runner, dir)
	if err != nil {
		return err
	}

	report := collab.AssessPull(entries, gitx.HasLFS(runner, dir))
	printWarnings(report.Warnings)
	if !report.CanPull {
		return fmt.Errorf("pull preflight failed")
	}

	if err := gitx.Pull(runner, dir); err != nil {
		return err
	}

	fmt.Println("Pull completed")
	return nil
}
