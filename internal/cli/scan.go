package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/scan"
	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "Inspect CAD assets and repository risks",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runScan(dir)
		},
	}
}

func runScan(dir string) error {
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

	files, usedMetadata, err := scan.LoadFiles(repoRoot, cfg.TrackedExtensions)
	if err != nil {
		return err
	}

	attributesData, err := os.ReadFile(filepath.Join(repoRoot, ".gitattributes"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	fmt.Print(formatScanReport(scan.BuildReport(files, string(attributesData), usedMetadata)))
	return nil
}

func formatScanReport(report scan.Report) string {
	var builder strings.Builder
	builder.WriteString("Repo Summary\n")
	builder.WriteString(fmt.Sprintf("- Total CAD files: %d\n", report.TotalFiles))
	builder.WriteString(fmt.Sprintf("- CAD types: %d\n", len(report.ByType)))
	builder.WriteString(fmt.Sprintf("- Locking recommended: %d\n", len(report.LockingRecommended)))
	builder.WriteString(fmt.Sprintf("- Git LFS expected: %d\n", len(report.GitLFSExpected)))
	builder.WriteString(fmt.Sprintf("- LFS warnings: %d\n", len(report.LFSWarnings)))
	if report.UsedMetadata {
		builder.WriteString("- Data source: metadata manifest\n")
	} else {
		builder.WriteString("- Data source: live scan\n")
	}

	if report.TotalFiles == 0 {
		builder.WriteString("No CAD files found for configured extensions\n")
		return builder.String()
	}

	builder.WriteString("Counts By Type\n")
	for _, count := range report.ByType {
		builder.WriteString(fmt.Sprintf("- %s: %d\n", count.TypeName, count.Count))
	}

	if len(report.TopDirectories) > 0 {
		builder.WriteString("Top Directories\n")
		for _, directory := range report.TopDirectories {
			builder.WriteString(fmt.Sprintf("- %s: %d\n", directory.Path, directory.Count))
		}
	}

	builder.WriteString("Locking Recommendations\n")
	if len(report.LockingRecommended) == 0 {
		builder.WriteString("- None\n")
	} else {
		for _, path := range report.DisplayedLocking {
			builder.WriteString(fmt.Sprintf("- %s\n", path))
		}
		if len(report.LockingRecommended) > len(report.DisplayedLocking) {
			builder.WriteString(fmt.Sprintf("- ... and %d more\n", len(report.LockingRecommended)-len(report.DisplayedLocking)))
		}
	}

	builder.WriteString("Git LFS Expected\n")
	if len(report.GitLFSExpected) == 0 {
		builder.WriteString("- None\n")
	} else {
		for _, path := range report.DisplayedLFSExpected {
			builder.WriteString(fmt.Sprintf("- %s\n", path))
		}
		if len(report.GitLFSExpected) > len(report.DisplayedLFSExpected) {
			builder.WriteString(fmt.Sprintf("- ... and %d more\n", len(report.GitLFSExpected)-len(report.DisplayedLFSExpected)))
		}
	}

	builder.WriteString("LFS Status\n")
	if len(report.LFSWarnings) == 0 {
		builder.WriteString("- All expected CAD files have matching .gitattributes LFS rules\n")
	} else {
		for _, warning := range report.LFSWarnings {
			builder.WriteString(fmt.Sprintf("- %s missing .gitattributes LFS rule for %s\n", warning.Path, warning.Extension))
		}
	}

	builder.WriteString("Largest CAD Files\n")
	for _, file := range report.LargestFiles {
		builder.WriteString(fmt.Sprintf("- %s | %s\n", file.Path, formatBytes(file.SizeBytes)))
	}

	return builder.String()
}

func formatBytes(size int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case size >= gb:
		return fmt.Sprintf("%.1f GB", float64(size)/float64(gb))
	case size >= mb:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(mb))
	case size >= kb:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(kb))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
