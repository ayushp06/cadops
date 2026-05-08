package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	commitcheck "github.com/cadops/cadops/internal/commit"
	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/metadata"
	"github.com/cadops/cadops/internal/preview"
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
	printCommitMetadataPreviewWarnings(repoRoot, entries)
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

func printCommitMetadataPreviewWarnings(root string, entries []gitx.StatusEntry) {
	cadPaths := statusEntryPaths(summarizeStatus(entries).CADFiles)
	if len(cadPaths) == 0 {
		return
	}

	if coverage, err := metadata.CheckCoverage(root, cadPaths); err == nil {
		if len(coverage.Missing) > 0 {
			fmt.Printf("Warning: Metadata: missing records for %s\n", strings.Join(coverage.Missing, ", "))
		}
		if len(coverage.Stale) > 0 {
			fmt.Printf("Warning: Metadata: stale records for %s\n", strings.Join(coverage.Stale, ", "))
		}
	}

	manifest, err := preview.Load(root)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Warning: Previews: missing records for %s\n", strings.Join(cadPaths, ", "))
		}
		return
	}

	missing := make([]string, 0)
	stale := make([]string, 0)
	for _, path := range cadPaths {
		record, ok := preview.Lookup(manifest, path)
		if !ok {
			missing = append(missing, path)
			continue
		}
		if preview.EffectiveStatus(root, record) == preview.StatusStale {
			stale = append(stale, path)
		}
	}
	if len(missing) > 0 {
		fmt.Printf("Warning: Previews: missing records for %s\n", strings.Join(missing, ", "))
	}
	if len(stale) > 0 {
		fmt.Printf("Warning: Previews: stale records for %s\n", strings.Join(stale, ", "))
	}
}
