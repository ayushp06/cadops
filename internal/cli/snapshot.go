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

type snapshotMetadataUpdater func(root string, extensions []string) (snapshot.MetadataUpdate, error)

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
	return runSnapshotWithMetadataUpdater(dir, now, snapshot.RefreshMetadata)
}

func runSnapshotWithMetadataUpdater(dir string, now time.Time, updateMetadata snapshotMetadataUpdater) error {
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

	plan, err := snapshot.BuildPlan(entries, cfg.TrackedExtensions, now)
	if err != nil {
		return err
	}

	for _, path := range plan.Paths {
		if err := gitx.StagePath(runner, repoRoot, path); err != nil {
			return fmt.Errorf("stage %s: %w", path, err)
		}
	}

	commitPaths := append([]string{}, plan.Paths...)
	metadataMessage := ""
	if updateMetadata != nil {
		update, err := updateMetadata(repoRoot, cfg.TrackedExtensions)
		if err != nil {
			fmt.Printf("Warning: metadata update failed: %v\n", err)
		} else {
			if err := gitx.StagePath(runner, repoRoot, update.Path); err != nil {
				fmt.Printf("Warning: metadata update failed: stage %s: %v\n", update.Path, err)
			} else {
				commitPaths = append(commitPaths, update.Path)
				metadataMessage = fmt.Sprintf("Metadata updated for %d CAD files\n", update.RecordCount)
			}
		}
	}

	if err := gitx.CommitPaths(runner, repoRoot, plan.Message, commitPaths); err != nil {
		return err
	}

	fmt.Printf("Created snapshot commit: %s\n", plan.Message)
	fmt.Printf("CAD files: %d\n", len(plan.Paths))
	if metadataMessage != "" {
		fmt.Print(metadataMessage)
	}
	return nil
}
