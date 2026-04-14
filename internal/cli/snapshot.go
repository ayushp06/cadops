package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/snapshot"
	"github.com/spf13/cobra"
)

func newSnapshotCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "snapshot",
		Short: "Create a CAD-focused snapshot commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runSnapshot(dir, time.Now())
		},
	}
}

func runSnapshot(dir string, now time.Time) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	cfg, err := config.Load(filepath.Join(dir, config.FileName))
	if err != nil {
		return err
	}

	entries, err := gitx.StatusPorcelain(runner, dir)
	if err != nil {
		return err
	}

	plan, err := snapshot.BuildPlan(entries, cfg.TrackedExtensions, now)
	if err != nil {
		return err
	}

	for _, path := range plan.Paths {
		if err := gitx.StagePath(runner, dir, path); err != nil {
			return fmt.Errorf("stage %s: %w", path, err)
		}
	}

	if err := gitx.CommitPaths(runner, dir, plan.Message, plan.Paths); err != nil {
		return err
	}

	fmt.Printf("Created snapshot commit: %s\n", plan.Message)
	fmt.Printf("CAD files: %d\n", len(plan.Paths))
	return nil
}
