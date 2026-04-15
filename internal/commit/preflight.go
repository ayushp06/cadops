package commit

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/collab"
	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
)

// Report captures CAD-aware commit readiness and warnings.
type Report struct {
	Warnings           []collab.Warning
	CanCommit          bool
	HasStagedChanges   bool
	HasUnstagedChanges bool
}

// Assess evaluates current repository state before a commit.
func Assess(cfg config.Config, entries []gitx.StatusEntry, attributes string, lockedPaths map[string]bool) Report {
	report := Report{CanCommit: true}

	stagedEntries, unstagedEntries := splitEntries(entries)
	report.HasStagedChanges = len(stagedEntries) > 0
	report.HasUnstagedChanges = len(unstagedEntries) > 0

	if !report.HasStagedChanges {
		report.CanCommit = false
		return report
	}

	if report.HasUnstagedChanges {
		report.Warnings = append(report.Warnings, collab.Warning{
			Title:   "Unstaged changes",
			Details: "local changes are present and will not be included in this commit",
		})
	}

	uncovered := FindUncoveredChangedCADFiles(entries, attributes)
	if len(uncovered) > 0 {
		report.Warnings = append(report.Warnings, collab.Warning{
			Title:   "LFS coverage",
			Details: "changed CAD files are missing matching Git LFS rules: " + strings.Join(uncovered, ", "),
		})
	}

	if cfg.LockingEnabled {
		missingLocks := FindRecommendedLockWarnings(entries, lockedPaths)
		if len(missingLocks) > 0 {
			report.Warnings = append(report.Warnings, collab.Warning{
				Title:   "Locking",
				Details: "recommended-lock CAD files are modified without a local Git LFS lock: " + strings.Join(missingLocks, ", "),
			})
		}
	}

	return report
}

// FindUncoveredChangedCADFiles returns changed CAD files without matching LFS rules.
func FindUncoveredChangedCADFiles(entries []gitx.StatusEntry, attributes string) []string {
	seenExtensions := make(map[string]bool)
	uncovered := make([]string, 0)

	for _, entry := range entries {
		path := changedCADPath(entry)
		if path == "" {
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

// FindRecommendedLockWarnings returns modified CAD files that recommend locking but are not locally locked.
func FindRecommendedLockWarnings(entries []gitx.StatusEntry, lockedPaths map[string]bool) []string {
	warnings := make([]string, 0)
	seen := make(map[string]bool)

	for _, entry := range entries {
		path := changedCADPath(entry)
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true

		fileType, ok := cad.Lookup(filepath.Ext(path))
		if !ok || !fileType.RecommendLocking {
			continue
		}
		if lockedPaths[path] {
			continue
		}
		warnings = append(warnings, path)
	}

	sort.Strings(warnings)
	return warnings
}

func splitEntries(entries []gitx.StatusEntry) ([]gitx.StatusEntry, []gitx.StatusEntry) {
	staged := make([]gitx.StatusEntry, 0, len(entries))
	unstaged := make([]gitx.StatusEntry, 0, len(entries))

	for _, entry := range entries {
		if isStagedEntry(entry) {
			staged = append(staged, entry)
		}
		if isUnstagedEntry(entry) {
			unstaged = append(unstaged, entry)
		}
	}

	return staged, unstaged
}

func changedCADPath(entry gitx.StatusEntry) string {
	if cad.IsCADPath(entry.Path) {
		return entry.Path
	}
	if cad.IsCADPath(entry.OldPath) {
		return entry.OldPath
	}
	return ""
}

func isStagedEntry(entry gitx.StatusEntry) bool {
	if len(entry.Code) < 2 {
		return false
	}
	return entry.Code[0] != ' ' && entry.Code[0] != '?'
}

func isUnstagedEntry(entry gitx.StatusEntry) bool {
	if len(entry.Code) < 2 {
		return false
	}
	return entry.Code[1] != ' '
}
