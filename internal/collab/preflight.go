package collab

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/gitx"
)

// Warning describes a non-fatal preflight issue.
type Warning struct {
	Title   string
	Details string
}

// PushReport captures the CAD-aware state before push.
type PushReport struct {
	Warnings []Warning
	CanPush  bool
}

// PullReport captures the CAD-aware state before pull.
type PullReport struct {
	Warnings []Warning
	CanPull  bool
}

// AssessPush builds a push preflight report from repository state.
func AssessPush(entries []gitx.StatusEntry, trackedFiles []string, attributes string, hasRemote bool) PushReport {
	report := PushReport{CanPush: true}

	if cadCount := countCADEntries(entries); cadCount > 0 {
		report.Warnings = append(report.Warnings, Warning{
			Title:   "Local CAD changes",
			Details: "unstaged or uncommitted CAD changes are present",
		})
	}

	uncovered := FindUncoveredCADFiles(trackedFiles, attributes)
	if len(uncovered) > 0 {
		report.Warnings = append(report.Warnings, Warning{
			Title:   "LFS coverage",
			Details: "CAD files are missing matching Git LFS rules: " + strings.Join(uncovered, ", "),
		})
	}

	if !hasRemote {
		report.Warnings = append(report.Warnings, Warning{
			Title:   "Remote",
			Details: "no git remote is configured",
		})
		report.CanPush = false
	}

	return report
}

// AssessPull builds a pull preflight report from repository state.
func AssessPull(entries []gitx.StatusEntry, hasLFS bool) PullReport {
	report := PullReport{CanPull: true}

	if len(entries) > 0 {
		report.Warnings = append(report.Warnings, Warning{
			Title:   "Dirty working tree",
			Details: "local changes are present and may complicate pull results",
		})
	}

	if cadCount := countCADEntries(entries); cadCount > 0 {
		report.Warnings = append(report.Warnings, Warning{
			Title:   "Modified CAD files",
			Details: "local CAD files are modified",
		})
	}

	if !hasLFS {
		report.Warnings = append(report.Warnings, Warning{
			Title:   "Git LFS",
			Details: "git lfs is not available",
		})
		report.CanPull = false
	}

	return report
}

// FindUncoveredCADFiles returns tracked CAD files without a matching attribute rule.
func FindUncoveredCADFiles(trackedFiles []string, attributes string) []string {
	seenExtensions := make(map[string]bool)
	uncovered := make([]string, 0)

	for _, path := range trackedFiles {
		if !cad.IsCADPath(path) {
			continue
		}
		ext := strings.ToLower(filepath.Ext(path))
		if seenExtensions[ext] {
			continue
		}
		seenExtensions[ext] = true
		if strings.Contains(attributes, gitx.AttributeLine(ext)) {
			continue
		}
		uncovered = append(uncovered, path)
	}

	sort.Strings(uncovered)
	return uncovered
}

func countCADEntries(entries []gitx.StatusEntry) int {
	count := 0
	for _, entry := range entries {
		if cad.IsCADPath(entry.Path) {
			count++
		}
	}
	return count
}
