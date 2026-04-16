package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/history"
	"github.com/cadops/cadops/internal/metadata"
	"github.com/spf13/cobra"
)

const defaultHistoryLimit = 10
const historyMetadataManifestPath = ".cadops/metadata/manifest.json"

func newHistoryCmd() *cobra.Command {
	limit := defaultHistoryLimit
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show recent CAD-aware commit history",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runHistory(dir, limit)
		},
	}
	cmd.Flags().IntVarP(&limit, "limit", "n", defaultHistoryLimit, "Number of recent commits to show")
	return cmd
}

func runHistory(dir string, limit int) error {
	if limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}
	if !gitx.IsInstalled("git") {
		return fmt.Errorf("git is not installed or not on PATH")
	}

	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}
	if !gitx.HasCommits(runner, dir) {
		fmt.Println("No commits found")
		return nil
	}

	repoRoot, err := gitx.RepoRoot(runner, dir)
	if err != nil {
		return err
	}

	out, err := gitx.RecentHistory(runner, repoRoot, limit)
	if err != nil {
		return err
	}

	entries := history.ParseLog(out)
	loadManifest := newHistoryManifestLoader(runner, repoRoot)
	resolveParent := func(commit string) (string, error) {
		return gitx.FirstParent(runner, repoRoot, commit)
	}

	fmt.Print(history.FormatReport(history.BuildReport(entries, loadManifest, resolveParent)))
	return nil
}

func newHistoryManifestLoader(runner gitx.Runner, repoRoot string) history.ManifestLoader {
	cache := make(map[string]metadata.Manifest)
	missing := make(map[string]bool)

	return func(revision string) (metadata.Manifest, error) {
		if manifest, ok := cache[revision]; ok {
			return manifest, nil
		}
		if missing[revision] {
			return metadata.Manifest{}, os.ErrNotExist
		}

		data, err := gitx.ReadFileAtRevision(runner, repoRoot, revision, filepath.ToSlash(historyMetadataManifestPath))
		if err != nil {
			if isMissingRevisionFileError(err) {
				missing[revision] = true
				return metadata.Manifest{}, os.ErrNotExist
			}
			return metadata.Manifest{}, err
		}

		manifest, err := metadata.Parse(data)
		if err != nil {
			return metadata.Manifest{}, err
		}
		cache[revision] = manifest
		return manifest, nil
	}
}

func isMissingRevisionFileError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "exists on disk, but not in") ||
		strings.Contains(msg, "path not in the working tree")
}
