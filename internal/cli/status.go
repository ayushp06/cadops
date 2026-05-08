package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/metadata"
	"github.com/cadops/cadops/internal/preview"
	"github.com/spf13/cobra"
)

type statusSummary struct {
	CADFiles        []gitx.StatusEntry
	NonCADFiles     []gitx.StatusEntry
	UncoveredCADLFS []string
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show CAD-aware repository status",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runStatus(dir)
		},
	}
}

func runStatus(dir string) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}
	repoRoot, err := gitx.RepoRoot(runner, dir)
	if err != nil {
		return err
	}

	entries, err := gitx.StatusPorcelain(runner, repoRoot)
	if err != nil {
		return err
	}
	summary := summarizeStatus(entries)

	if len(entries) == 0 {
		fmt.Println("Working tree clean")
	} else {
		fmt.Printf("Changed files: %d\n", len(entries))
		printStatusGroup("CAD files", summary.CADFiles)
		printStatusGroup("Non-CAD files", summary.NonCADFiles)
	}

	cadPaths := statusEntryPaths(summary.CADFiles)
	printMetadataStatus(repoRoot, cadPaths)
	printPreviewStatus(repoRoot, cadPaths)

	uncovered, err := findUncoveredLFS(repoRoot, summary.CADFiles)
	if err != nil {
		return err
	}
	if len(uncovered) > 0 {
		fmt.Println("Warnings:")
		for _, path := range uncovered {
			fmt.Printf("- %s is a CAD file without a matching .gitattributes LFS rule\n", path)
		}
	}

	return nil
}

func statusEntryPaths(entries []gitx.StatusEntry) []string {
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Path != "" {
			paths = append(paths, entry.Path)
		}
	}
	sort.Strings(paths)
	return paths
}

func printMetadataStatus(root string, cadPaths []string) {
	if len(cadPaths) == 0 {
		return
	}
	coverage, err := metadata.CheckCoverage(root, cadPaths)
	if err != nil {
		fmt.Printf("Metadata: unavailable (%v)\n", err)
		return
	}
	fmt.Printf("Metadata: missing %d, stale %d\n", len(coverage.Missing), len(coverage.Stale))
}

func printPreviewStatus(root string, cadPaths []string) {
	if len(cadPaths) == 0 {
		return
	}
	manifest, err := preview.Load(root)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Previews: missing %d, stale 0\n", len(cadPaths))
			return
		}
		fmt.Printf("Previews: unavailable (%v)\n", err)
		return
	}

	missing := 0
	stale := 0
	for _, path := range cadPaths {
		record, ok := preview.Lookup(manifest, path)
		if !ok {
			missing++
			continue
		}
		if preview.EffectiveStatus(root, record) == preview.StatusStale {
			stale++
		}
	}
	fmt.Printf("Previews: missing %d, stale %d\n", missing, stale)
}

func summarizeStatus(entries []gitx.StatusEntry) statusSummary {
	summary := statusSummary{
		CADFiles:    make([]gitx.StatusEntry, 0),
		NonCADFiles: make([]gitx.StatusEntry, 0),
	}

	for _, entry := range entries {
		if cad.IsCADPath(entry.Path) {
			summary.CADFiles = append(summary.CADFiles, entry)
		} else {
			summary.NonCADFiles = append(summary.NonCADFiles, entry)
		}
	}
	return summary
}

func printStatusGroup(title string, entries []gitx.StatusEntry) {
	fmt.Printf("%s: %d\n", title, len(entries))
	for _, entry := range entries {
		fmt.Printf("- [%s] %s\n", strings.TrimSpace(entry.Code), entry.Path)
	}
}

func findUncoveredLFS(dir string, cadEntries []gitx.StatusEntry) ([]string, error) {
	attributesData, err := os.ReadFile(filepath.Join(dir, ".gitattributes"))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	attributes := string(attributesData)

	var uncovered []string
	seen := make(map[string]bool, len(cadEntries))
	for _, entry := range cadEntries {
		ext := strings.ToLower(filepath.Ext(entry.Path))
		if seen[ext] {
			continue
		}
		seen[ext] = true
		if !strings.Contains(attributes, gitx.AttributeLine(ext)) {
			uncovered = append(uncovered, entry.Path)
		}
	}
	sort.Strings(uncovered)
	return uncovered, nil
}
