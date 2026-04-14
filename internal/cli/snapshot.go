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

	paths := snapshot.SelectPaths(entries, cfg.TrackedExtensions)
	if len(paths) == 0 {
		return fmt.Errorf("no changed CAD files to snapshot")
	}

	for _, path := range paths {
		if err := gitx.StagePath(runner, dir, path); err != nil {
			return fmt.Errorf("stage %s: %w", path, err)
		}
	}

	message := snapshot.BuildMessage(now)
	if err := gitx.CommitPaths(runner, dir, message, paths); err != nil {
		return err
	}

	fmt.Printf("Created snapshot commit: %s\n", message)
	fmt.Printf("CAD files: %d\n", len(paths))
	return nil
}
