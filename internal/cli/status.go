package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/gitx"
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

	entries, err := gitx.StatusPorcelain(runner, dir)
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

	uncovered, err := findUncoveredLFS(dir, summary.CADFiles)
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
