package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cadops/cadops/internal/config"
	caddiff "github.com/cadops/cadops/internal/diff"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Show a CAD-aware summary of repository changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runDiff(dir)
		},
	}
}

func runDiff(dir string) error {
	if !gitx.IsInstalled("git") {
		return fmt.Errorf("git is not installed or not on PATH")
	}

	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	repoRoot, err := gitx.RepoRoot(runner, dir)
	if err != nil {
		return err
	}

	cfg, err := config.Load(filepath.Join(repoRoot, config.FileName))
	if err != nil {
		return err
	}

	statusEntries, err := gitx.StatusPorcelain(runner, repoRoot)
	if err != nil {
		return err
	}

	entries := caddiff.BuildEntries(statusEntries, cfg.TrackedExtensions)
	fmt.Print(caddiff.FormatReport(caddiff.BuildReport(repoRoot, entries)))
	return nil
}
