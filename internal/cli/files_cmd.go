package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/config"
	cadfiles "github.com/cadops/cadops/internal/files"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/spf13/cobra"
)

func newFilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "files",
		Short: "List CAD-relevant repository files",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runFiles(dir)
		},
	}
}

func runFiles(dir string) error {
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

	entries, err := cadfiles.Scan(repoRoot, cfg.TrackedExtensions)
	if err != nil {
		return err
	}

	fmt.Print(formatFilesReport(cadfiles.GroupEntries(entries)))
	return nil
}

func formatFilesReport(groups []cadfiles.Group) string {
	if len(groups) == 0 {
		return "No CAD files found for configured extensions\n"
	}

	var builder strings.Builder
	total := 0
	for _, group := range groups {
		total += len(group.Entries)
	}

	builder.WriteString(fmt.Sprintf("CAD files: %d\n", total))
	for _, group := range groups {
		builder.WriteString(group.TypeName)
		builder.WriteString("\n")
		for _, entry := range group.Entries {
			locking := "no"
			if entry.RecommendLocking {
				locking = "yes"
			}
			builder.WriteString(fmt.Sprintf("- %s | type: %s | lock: %s\n", entry.Path, entry.TypeName, locking))
		}
	}

	return builder.String()
}
