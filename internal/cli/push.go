package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cadops/cadops/internal/collab"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/spf13/cobra"
)

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Run CAD-aware checks before git push",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runPush(dir)
		},
	}
}

func runPush(dir string) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	entries, err := gitx.StatusPorcelain(runner, dir)
	if err != nil {
		return err
	}
	trackedFiles, err := gitx.ListTrackedFiles(runner, dir)
	if err != nil {
		return err
	}
	attributesData, err := os.ReadFile(filepath.Join(dir, ".gitattributes"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	report := collab.AssessPush(entries, trackedFiles, string(attributesData), gitx.HasRemote(runner, dir))
	printWarnings(report.Warnings)
	if !report.CanPush {
		return fmt.Errorf("push preflight failed")
	}

	if err := gitx.Push(runner, dir); err != nil {
		return err
	}

	fmt.Println("Push completed")
	return nil
}
