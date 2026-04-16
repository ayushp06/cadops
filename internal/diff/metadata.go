package diff

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/metadata"
)

// Report is a metadata-aware diff summary for terminal rendering.
type Report struct {
	CAD      []DetailedEntry
	Other    []Entry
	Warnings []string
}

// DetailedEntry combines a diff entry with optional metadata context.
type DetailedEntry struct {
	Entry    Entry
	Metadata MetadataDetails
}

// MetadataDetails stores available metadata context for a changed CAD file.
type MetadataDetails struct {
	Current            metadata.Record
	HasCurrent         bool
	Previous           metadata.Record
	HasPrevious        bool
	Comparison         Comparison
	PreviousLookupPath string
}

// Comparison captures compact change signals derived from two metadata records.
type Comparison struct {
	ChecksumChanged bool
	SizeDeltaBytes  int64
	HasSizeDelta    bool
}

// BuildReport groups changed entries and enriches CAD files with metadata when
// the current manifest and filesystem provide enough information.
func BuildReport(root string, entries []Entry) Report {
	summary := Summarize(entries)
	report := Report{
		CAD:   make([]DetailedEntry, 0, len(summary.CAD)),
		Other: summary.Other,
	}

	manifest, hasManifest, manifestWarnings := loadManifest(root)
	report.Warnings = append(report.Warnings, manifestWarnings...)

	for _, entry := range summary.CAD {
		detailed, warnings := enrichEntry(root, manifest, hasManifest, entry)
		report.CAD = append(report.CAD, detailed)
		report.Warnings = append(report.Warnings, warnings...)
	}

	return report
}

// CompareRecords returns compact comparison facts for previous versus current
// metadata records.
func CompareRecords(previous, current metadata.Record) Comparison {
	return Comparison{
		ChecksumChanged: previous.SHA256 != "" && current.SHA256 != "" && previous.SHA256 != current.SHA256,
		SizeDeltaBytes:  current.SizeBytes - previous.SizeBytes,
		HasSizeDelta:    true,
	}
}

func enrichEntry(root string, manifest metadata.Manifest, hasManifest bool, entry Entry) (DetailedEntry, []string) {
	detailed := DetailedEntry{Entry: entry}
	warnings := make([]string, 0, 1)

	if hasManifest {
		previous, lookupPath, ok := lookupPreviousRecord(manifest, entry)
		if ok {
			detailed.Metadata.Previous = previous
			detailed.Metadata.HasPrevious = true
			detailed.Metadata.PreviousLookupPath = lookupPath
		}
	}

	if entry.Kind != KindDeleted {
		current, err := buildCurrentRecord(root, entry.Path)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("metadata unavailable for %s", entry.Path))
			return detailed, warnings
		}
		detailed.Metadata.Current = current
		detailed.Metadata.HasCurrent = true
	}

	if detailed.Metadata.HasPrevious && detailed.Metadata.HasCurrent {
		detailed.Metadata.Comparison = CompareRecords(detailed.Metadata.Previous, detailed.Metadata.Current)
	}

	return detailed, warnings
}

func loadManifest(root string) (metadata.Manifest, bool, []string) {
	manifest, err := metadata.Load(root)
	if err == nil {
		return manifest, true, nil
	}
	if os.IsNotExist(err) {
		return metadata.Manifest{}, false, nil
	}
	return metadata.Manifest{}, false, []string{"metadata store unavailable; showing standard diff context only"}
}

func lookupPreviousRecord(manifest metadata.Manifest, entry Entry) (metadata.Record, string, bool) {
	paths := []string{entry.Path}
	if entry.OldPath != "" && entry.OldPath != entry.Path {
		paths = append(paths, entry.OldPath)
	}

	for _, path := range paths {
		record, ok := metadata.Lookup(manifest, path)
		if ok {
			return record, path, true
		}
	}

	return metadata.Record{}, "", false
}

func buildCurrentRecord(root, relPath string) (metadata.Record, error) {
	absPath := filepath.Join(root, filepath.FromSlash(relPath))
	info, err := os.Stat(absPath)
	if err != nil {
		return metadata.Record{}, err
	}
	return metadata.BuildRecord(root, relPath, info)
}

// FormatReport renders a compact diff report with metadata-aware CAD details.
func FormatReport(report Report) string {
	total := len(report.CAD) + len(report.Other)
	if total == 0 {
		return "No repository changes\n"
	}

	var builder strings.Builder
	writeCADGroup(&builder, report.CAD)
	writeDiffGroup(&builder, "Other changes", report.Other)
	if len(report.Warnings) > 0 {
		builder.WriteString("Warnings:\n")
		for _, warning := range report.Warnings {
			builder.WriteString("  - ")
			builder.WriteString(warning)
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func writeDiffGroup(builder *strings.Builder, title string, entries []Entry) {
	if len(entries) == 0 {
		return
	}
	builder.WriteString(title)
	builder.WriteString(":\n")
	for _, entry := range entries {
		builder.WriteString(fmt.Sprintf("  %s %s\n", entry.Kind, DisplayPath(entry)))
	}
}

func writeCADGroup(builder *strings.Builder, entries []DetailedEntry) {
	if len(entries) == 0 {
		return
	}
	builder.WriteString("CAD changes:\n")
	for _, entry := range entries {
		builder.WriteString("  ")
		builder.WriteString(string(entry.Entry.Kind))
		builder.WriteString(" ")
		builder.WriteString(DisplayPath(entry.Entry))
		if details := formatMetadataDetails(entry.Metadata); details != "" {
			builder.WriteString(" [")
			builder.WriteString(details)
			builder.WriteString("]")
		}
		builder.WriteString("\n")
	}
}

func formatMetadataDetails(details MetadataDetails) string {
	record, ok := preferredRecord(details)
	if !ok {
		return ""
	}

	parts := []string{
		record.TypeName,
		"lock " + yesNo(record.LockingRecommended),
		"LFS " + yesNo(record.GitLFSExpected),
	}

	if details.HasPrevious && details.HasCurrent && details.Comparison.ChecksumChanged {
		parts = append(parts, "checksum changed")
	}
	if details.HasPrevious && details.HasCurrent && details.Comparison.HasSizeDelta && details.Comparison.SizeDeltaBytes != 0 {
		parts = append(parts, "size "+formatSizeDelta(details.Comparison.SizeDeltaBytes))
	}
	if !details.HasCurrent && details.HasPrevious {
		parts = append(parts, "stored metadata")
	}

	return strings.Join(parts, "; ")
}

func preferredRecord(details MetadataDetails) (metadata.Record, bool) {
	if details.HasCurrent {
		return details.Current, true
	}
	if details.HasPrevious {
		return details.Previous, true
	}
	return metadata.Record{}, false
}

// FormatSizeDelta renders a signed byte delta.
func FormatSizeDelta(delta int64) string {
	return formatSizeDelta(delta)
}

func formatSizeDelta(delta int64) string {
	sign := "+"
	if delta < 0 {
		sign = "-"
		delta = -delta
	}
	return fmt.Sprintf("%s%d B", sign, delta)
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
