package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	fmt.Print(formatDiffReport(caddiff.Summarize(entries)))
	return nil
}

func formatDiffReport(summary caddiff.Summary) string {
	total := len(summary.CAD) + len(summary.Other)
	if total == 0 {
		return "No repository changes\n"
	}

	var builder strings.Builder
	writeDiffGroup(&builder, "CAD changes", summary.CAD)
	writeDiffGroup(&builder, "Other changes", summary.Other)
	return builder.String()
}

func writeDiffGroup(builder *strings.Builder, title string, entries []caddiff.Entry) {
	if len(entries) == 0 {
		return
	}
	builder.WriteString(title)
	builder.WriteString(":\n")
	for _, entry := range entries {
		builder.WriteString(fmt.Sprintf("  %s %s\n", entry.Kind, caddiff.DisplayPath(entry)))
	}
}
