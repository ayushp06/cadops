package history

import (
	"fmt"
	"os"
	"strings"

	"github.com/cadops/cadops/internal/metadata"
)

// ManifestLoader loads a metadata manifest for a revision.
type ManifestLoader func(revision string) (metadata.Manifest, error)

// ParentResolver returns the first parent hash for a commit.
type ParentResolver func(commit string) (string, error)

// Report is a metadata-aware history view.
type Report struct {
	Entries  []DetailedEntry
	Warnings []string
}

// DetailedEntry combines a commit entry with metadata-aware CAD file details.
type DetailedEntry struct {
	Entry    Entry
	CADFiles []CADFileDetail
}

// CADFileDetail describes one changed CAD file plus optional metadata context.
type CADFileDetail struct {
	Path     string
	Metadata FileMetadata
}

// FileMetadata captures history-safe metadata context for a CAD file.
type FileMetadata struct {
	Current           metadata.Record
	HasCurrent        bool
	Previous          metadata.Record
	HasPrevious       bool
	MetadataAvailable bool
	ChecksumChanged   bool
	SizeDeltaBytes    int64
	HasSizeDelta      bool
}

// BuildReport enriches parsed history entries with commit-scoped metadata when
// manifests are available in Git history.
func BuildReport(entries []Entry, loadManifest ManifestLoader, resolveParent ParentResolver) Report {
	report := Report{
		Entries:  make([]DetailedEntry, 0, len(entries)),
		Warnings: make([]string, 0),
	}
	if len(entries) == 0 {
		return report
	}

	for _, entry := range entries {
		detailed := DetailedEntry{
			Entry:    entry,
			CADFiles: make([]CADFileDetail, 0, len(entry.CADFiles)),
		}

		currentManifest, hasCurrentManifest, currentWarning := safeLoadManifest(loadManifest, entry.Hash)
		if currentWarning != "" {
			report.Warnings = append(report.Warnings, currentWarning)
		}

		parentHash, parentWarning := safeResolveParent(resolveParent, entry.Hash)
		if parentWarning != "" {
			report.Warnings = append(report.Warnings, parentWarning)
		}

		var previousManifest metadata.Manifest
		hasPreviousManifest := false
		if parentHash != "" {
			var previousWarning string
			previousManifest, hasPreviousManifest, previousWarning = safeLoadManifest(loadManifest, parentHash)
			if previousWarning != "" {
				report.Warnings = append(report.Warnings, previousWarning)
			}
		}

		for _, path := range entry.CADFiles {
			file := CADFileDetail{Path: path}

			if hasCurrentManifest {
				if record, ok := metadata.Lookup(currentManifest, path); ok {
					file.Metadata.Current = record
					file.Metadata.HasCurrent = true
					file.Metadata.MetadataAvailable = true
				}
			}
			if hasPreviousManifest {
				if record, ok := metadata.Lookup(previousManifest, path); ok {
					file.Metadata.Previous = record
					file.Metadata.HasPrevious = true
				}
			}
			if file.Metadata.HasCurrent && file.Metadata.HasPrevious {
				file.Metadata.ChecksumChanged = file.Metadata.Current.SHA256 != "" &&
					file.Metadata.Previous.SHA256 != "" &&
					file.Metadata.Current.SHA256 != file.Metadata.Previous.SHA256
				file.Metadata.SizeDeltaBytes = file.Metadata.Current.SizeBytes - file.Metadata.Previous.SizeBytes
				file.Metadata.HasSizeDelta = true
			}
			if !file.Metadata.MetadataAvailable && file.Metadata.HasPrevious {
				file.Metadata.Current = file.Metadata.Previous
				file.Metadata.HasCurrent = true
				file.Metadata.MetadataAvailable = true
			}

			detailed.CADFiles = append(detailed.CADFiles, file)
		}

		report.Entries = append(report.Entries, detailed)
	}

	return report
}

// FormatReport renders a compact terminal-friendly history view.
func FormatReport(report Report) string {
	if len(report.Entries) == 0 {
		return "No commits found\n"
	}

	var builder strings.Builder
	for i, entry := range report.Entries {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(fmt.Sprintf("%s  %s  %s\n", entry.Entry.ShortHash, entry.Entry.Date, entry.Entry.Message))
		builder.WriteString("  CAD files:\n")
		if len(entry.CADFiles) == 0 {
			builder.WriteString("    - none\n")
			continue
		}
		for _, file := range entry.CADFiles {
			builder.WriteString("    - ")
			builder.WriteString(file.Path)
			if details := formatFileMetadata(file.Metadata); details != "" {
				builder.WriteString(" [")
				builder.WriteString(details)
				builder.WriteString("]")
			}
			builder.WriteString("\n")
		}
	}
	if len(report.Warnings) > 0 {
		builder.WriteString("\nWarnings:\n")
		for _, warning := range report.Warnings {
			builder.WriteString("  - ")
			builder.WriteString(warning)
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func safeResolveParent(resolveParent ParentResolver, commit string) (string, string) {
	if resolveParent == nil {
		return "", ""
	}
	parent, err := resolveParent(commit)
	if err != nil {
		return "", fmt.Sprintf("metadata lookup failed for %s", shortHash(commit))
	}
	return parent, ""
}

func safeLoadManifest(loadManifest ManifestLoader, revision string) (metadata.Manifest, bool, string) {
	if loadManifest == nil || revision == "" {
		return metadata.Manifest{}, false, ""
	}
	manifest, err := loadManifest(revision)
	if err == nil {
		return manifest, true, ""
	}
	if os.IsNotExist(err) {
		return metadata.Manifest{}, false, ""
	}
	return metadata.Manifest{}, false, fmt.Sprintf("metadata lookup failed for %s", shortHash(revision))
}

func formatFileMetadata(details FileMetadata) string {
	if !details.MetadataAvailable {
		return "metadata unavailable"
	}

	parts := []string{
		details.Current.TypeName,
		"size " + formatSize(details.Current.SizeBytes),
	}
	if details.HasPrevious && details.ChecksumChanged {
		parts = append(parts, "checksum changed")
	}
	if details.HasPrevious && details.HasSizeDelta && details.SizeDeltaBytes != 0 {
		parts = append(parts, "delta "+formatSizeDelta(details.SizeDeltaBytes))
	}
	return strings.Join(parts, "; ")
}

func formatSize(sizeBytes int64) string {
	return fmt.Sprintf("%d B", sizeBytes)
}

func formatSizeDelta(delta int64) string {
	sign := "+"
	if delta < 0 {
		sign = "-"
		delta = -delta
	}
	return fmt.Sprintf("%s%d B", sign, delta)
}
