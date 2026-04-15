package cli

import (
	"fmt"
	"os"
	"path/filepath"

	commitcheck "github.com/cadops/cadops/internal/commit"
	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/spf13/cobra"
)

func newCommitCmd() *cobra.Command {
	var message string

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Run CAD-aware checks before git commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runCommit(dir, message)
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
	_ = cmd.MarkFlagRequired("message")
	return cmd
}

func runCommit(dir, message string) error {
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

	entries, err := gitx.StatusPorcelain(runner, repoRoot)
	if err != nil {
		return err
	}

	attributesData, err := os.ReadFile(filepath.Join(repoRoot, ".gitattributes"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lockedPaths := map[string]bool{}
	if cfg.LockingEnabled && gitx.HasLFS(runner, repoRoot) {
		lockedPaths, err = gitx.ListLocalLocks(runner, repoRoot)
		if err != nil {
			lockedPaths = map[string]bool{}
		}
	}

	report := commitcheck.Assess(cfg, entries, string(attributesData), lockedPaths)
	printWarnings(report.Warnings)
	if !report.CanCommit {
		if report.HasUnstagedChanges {
			return fmt.Errorf("nothing staged to commit")
		}
		return fmt.Errorf("nothing to commit")
	}

	if err := gitx.Commit(runner, repoRoot, message); err != nil {
		return err
	}

	fmt.Println("Commit completed")
	return nil
}
