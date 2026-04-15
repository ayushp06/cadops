package scan

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/metadata"
	"github.com/cadops/cadops/internal/watch"
)

const (
	topDirectoryLimit = 5
	largestFileLimit  = 5
	pathListLimit     = 10
)

// File describes a scanned CAD file for repository inspection.
type File struct {
	Path               string
	Directory          string
	TypeName           string
	Extension          string
	SizeBytes          int64
	GitLFSExpected     bool
	LockingRecommended bool
}

// TypeCount summarizes files by CAD type.
type TypeCount struct {
	TypeName string
	Count    int
}

// DirectoryCount summarizes the most common CAD directories.
type DirectoryCount struct {
	Path  string
	Count int
}

// SizedFile summarizes a file with size for ranking.
type SizedFile struct {
	Path      string
	SizeBytes int64
}

// LFSWarning highlights a file missing expected LFS configuration.
type LFSWarning struct {
	Path      string
	Extension string
}

// Report is the aggregated repository inspection result.
type Report struct {
	TotalFiles           int
	ByType               []TypeCount
	LockingRecommended   []string
	GitLFSExpected       []string
	LFSWarnings          []LFSWarning
	LargestFiles         []SizedFile
	TopDirectories       []DirectoryCount
	UsedMetadata         bool
	DisplayedLocking     []string
	DisplayedLFSExpected []string
}

// LoadFiles loads scan input either from the metadata manifest when available
// or from a live repository walk.
func LoadFiles(root string, extensions []string) ([]File, bool, error) {
	manifest, err := metadata.Load(root)
	if err == nil {
		files := make([]File, 0, len(manifest.Records))
		for _, record := range manifest.Records {
			files = append(files, fileFromRecord(record))
		}
		sortFiles(files)
		return files, true, nil
	}
	if !os.IsNotExist(err) {
		return nil, false, err
	}

	files, err := Scan(root, extensions)
	if err != nil {
		return nil, false, err
	}
	return files, false, nil
}

// Scan walks the repository and builds lightweight file facts for scan.
func Scan(root string, extensions []string) ([]File, error) {
	filter := watch.NewFilter(extensions)
	files := make([]File, 0)

	err := filepath.WalkDir(root, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := dirEntry.Name()
		if dirEntry.IsDir() {
			if name == ".git" || path == filepath.Join(root, ".cadops") {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		if !filter.Match(relPath) {
			return nil
		}

		info, err := dirEntry.Info()
		if err != nil {
			return err
		}

		files = append(files, BuildFile(relPath, info.Size()))
		return nil
	})
	if err != nil {
		return nil, err
	}

	sortFiles(files)
	return files, nil
}

// BuildFile classifies a repository-relative path into a scan file fact.
func BuildFile(relPath string, sizeBytes int64) File {
	normalizedPath := filepath.ToSlash(filepath.Clean(relPath))
	extension := strings.ToLower(filepath.Ext(normalizedPath))
	fileType, ok := cad.Lookup(extension)
	file := File{
		Path:      normalizedPath,
		Directory: directoryForPath(normalizedPath),
		TypeName:  unknownTypeName(extension),
		Extension: extension,
		SizeBytes: sizeBytes,
	}
	if ok {
		file.TypeName = fileType.Name
		file.GitLFSExpected = fileType.UseLFS
		file.LockingRecommended = fileType.RecommendLocking
	}
	return file
}

// BuildReport aggregates a repository scan plus optional LFS attribute data.
func BuildReport(files []File, attributes string, usedMetadata bool) Report {
	report := Report{
		TotalFiles:         len(files),
		ByType:             CountByType(files),
		LockingRecommended: PathsMatching(files, func(file File) bool { return file.LockingRecommended }),
		GitLFSExpected:     PathsMatching(files, func(file File) bool { return file.GitLFSExpected }),
		LFSWarnings:        FindLFSWarnings(files, attributes),
		LargestFiles:       LargestFiles(files, largestFileLimit),
		TopDirectories:     TopDirectories(files, topDirectoryLimit),
		UsedMetadata:       usedMetadata,
	}
	report.DisplayedLocking = LimitPaths(report.LockingRecommended, pathListLimit)
	report.DisplayedLFSExpected = LimitPaths(report.GitLFSExpected, pathListLimit)
	return report
}

// CountByType returns stable type counts.
func CountByType(files []File) []TypeCount {
	counts := make(map[string]int, len(files))
	for _, file := range files {
		counts[file.TypeName]++
	}

	out := make([]TypeCount, 0, len(counts))
	for typeName, count := range counts {
		out = append(out, TypeCount{TypeName: typeName, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].TypeName < out[j].TypeName
	})
	return out
}

// TopDirectories returns the most common directories in stable rank order.
func TopDirectories(files []File, limit int) []DirectoryCount {
	counts := make(map[string]int)
	for _, file := range files {
		counts[file.Directory]++
	}

	out := make([]DirectoryCount, 0, len(counts))
	for path, count := range counts {
		out = append(out, DirectoryCount{Path: path, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Path < out[j].Path
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// LargestFiles returns the largest files in descending size order.
func LargestFiles(files []File, limit int) []SizedFile {
	out := make([]SizedFile, 0, len(files))
	for _, file := range files {
		out = append(out, SizedFile{
			Path:      file.Path,
			SizeBytes: file.SizeBytes,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].SizeBytes != out[j].SizeBytes {
			return out[i].SizeBytes > out[j].SizeBytes
		}
		return out[i].Path < out[j].Path
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// FindLFSWarnings returns files expected to use LFS but missing a matching
// canonical .gitattributes rule.
func FindLFSWarnings(files []File, attributes string) []LFSWarning {
	warnings := make([]LFSWarning, 0)
	for _, file := range files {
		if !file.GitLFSExpected {
			continue
		}
		if HasLFSRule(attributes, file.Extension) {
			continue
		}
		warnings = append(warnings, LFSWarning{
			Path:      file.Path,
			Extension: file.Extension,
		})
	}
	sort.Slice(warnings, func(i, j int) bool {
		return warnings[i].Path < warnings[j].Path
	})
	return warnings
}

// HasLFSRule reports whether .gitattributes contains the canonical rule for
// the given extension.
func HasLFSRule(attributes, extension string) bool {
	return strings.Contains(attributes, gitx.AttributeLine(extension))
}

// PathsMatching selects file paths using a stable filter.
func PathsMatching(files []File, keep func(File) bool) []string {
	out := make([]string, 0)
	for _, file := range files {
		if keep(file) {
			out = append(out, file.Path)
		}
	}
	sort.Strings(out)
	return out
}

// LimitPaths returns a stable prefix of the paths slice.
func LimitPaths(paths []string, limit int) []string {
	if limit <= 0 || len(paths) <= limit {
		out := make([]string, len(paths))
		copy(out, paths)
		return out
	}
	out := make([]string, limit)
	copy(out, paths[:limit])
	return out
}

func fileFromRecord(record metadata.Record) File {
	return File{
		Path:               record.Path,
		Directory:          directoryForPath(record.Path),
		TypeName:           record.TypeName,
		Extension:          record.Extension,
		SizeBytes:          record.SizeBytes,
		GitLFSExpected:     record.GitLFSExpected,
		LockingRecommended: record.LockingRecommended,
	}
}

func sortFiles(files []File) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
}

func directoryForPath(path string) string {
	dir := filepath.ToSlash(filepath.Dir(path))
	if dir == "." || dir == "" {
		return "."
	}
	return dir
}

func unknownTypeName(extension string) string {
	if extension == "" {
		return "Unknown CAD Type"
	}
	return "Unknown CAD Type (" + extension + ")"
}
