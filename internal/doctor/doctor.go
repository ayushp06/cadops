package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
)

// CheckLevel indicates the outcome severity.
type CheckLevel string

const (
	LevelPass CheckLevel = "PASS"
	LevelWarn CheckLevel = "WARN"
	LevelFail CheckLevel = "FAIL"
)

// CheckResult is a single doctor check outcome.
type CheckResult struct {
	Level   CheckLevel
	Name    string
	Details string
}

// Report contains all doctor checks and a convenience success state.
type Report struct {
	Results []CheckResult
}

// HasFailures reports whether any checks failed.
func (r Report) HasFailures() bool {
	for _, result := range r.Results {
		if result.Level == LevelFail {
			return true
		}
	}
	return false
}

// Run executes repository health checks.
func Run(dir string, runner gitx.Runner) Report {
	results := make([]CheckResult, 0, 8)

	if gitx.IsInstalled("git") {
		results = append(results, CheckResult{Level: LevelPass, Name: "Git", Details: "git is installed"})
	} else {
		results = append(results, CheckResult{Level: LevelFail, Name: "Git", Details: "git is not installed"})
	}

	if gitx.IsInstalled("git") {
		if _, err := runner.Run(dir, "git", "lfs", "version"); err == nil {
			results = append(results, CheckResult{Level: LevelPass, Name: "Git LFS", Details: "git lfs is installed"})
		} else {
			results = append(results, CheckResult{Level: LevelFail, Name: "Git LFS", Details: "git lfs is not installed"})
		}
	}

	if gitx.IsRepo(runner, dir) {
		results = append(results, CheckResult{Level: LevelPass, Name: "Repository", Details: "git repository detected"})
	} else {
		return Report{Results: append(results, CheckResult{Level: LevelFail, Name: "Repository", Details: "not a git repository"})}
	}

	configPath := filepath.Join(dir, config.FileName)
	if _, err := os.Stat(configPath); err == nil {
		results = append(results, CheckResult{Level: LevelPass, Name: "Config", Details: config.FileName + " exists"})
	} else {
		results = append(results, CheckResult{Level: LevelFail, Name: "Config", Details: config.FileName + " is missing"})
	}

	attributesPath := filepath.Join(dir, ".gitattributes")
	attributesData, err := os.ReadFile(attributesPath)
	if err != nil {
		results = append(results, CheckResult{Level: LevelFail, Name: "Attributes", Details: ".gitattributes is missing"})
	} else {
		missing := missingAttributeExtensions(string(attributesData), cad.SupportedExtensions())
		if len(missing) == 0 {
			results = append(results, CheckResult{Level: LevelPass, Name: "Attributes", Details: "all CAD extensions have LFS attributes"})
		} else {
			results = append(results, CheckResult{Level: LevelFail, Name: "Attributes", Details: "missing entries for " + strings.Join(missing, ", ")})
		}
	}

	untracked, err := FindUntrackedCADFiles(dir, runner)
	if err != nil {
		results = append(results, CheckResult{Level: LevelWarn, Name: "CAD Coverage", Details: "unable to inspect CAD files: " + err.Error()})
	} else if len(untracked) == 0 {
		results = append(results, CheckResult{Level: LevelPass, Name: "CAD Coverage", Details: "all CAD files are tracked by git"})
	} else {
		results = append(results, CheckResult{Level: LevelWarn, Name: "CAD Coverage", Details: "untracked CAD files: " + strings.Join(untracked, ", ")})
	}

	if gitx.HasRemote(runner, dir) {
		results = append(results, CheckResult{Level: LevelPass, Name: "Remote", Details: "git remote configured"})
	} else {
		results = append(results, CheckResult{Level: LevelWarn, Name: "Remote", Details: "no git remote configured"})
	}

	return Report{Results: results}
}

func missingAttributeExtensions(attributes string, extensions []string) []string {
	missing := make([]string, 0)
	for _, extension := range extensions {
		if !strings.Contains(attributes, gitx.AttributeLine(extension)) {
			missing = append(missing, extension)
		}
	}
	return missing
}

// FindUntrackedCADFiles discovers CAD files in the worktree that are not tracked by git.
func FindUntrackedCADFiles(dir string, runner gitx.Runner) ([]string, error) {
	tracked, err := gitx.ListTrackedFiles(runner, dir)
	if err != nil {
		return nil, fmt.Errorf("list tracked files: %w", err)
	}

	trackedSet := make(map[string]bool, len(tracked))
	for _, path := range tracked {
		trackedSet[path] = true
	}

	var untracked []string
	err = filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !cad.IsCADPath(rel) {
			return nil
		}
		if !trackedSet[rel] {
			untracked = append(untracked, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(untracked)
	return untracked, nil
}
