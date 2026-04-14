package locking

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/cad"
	"github.com/cadops/cadops/internal/gitx"
)

// Assessment describes a validated lock target and any setup warnings.
type Assessment struct {
	Path              string
	LockRecommended   bool
	ConfigurationWarn []string
}

// ResolveTarget validates that the path exists and returns a repo-relative path.
func ResolveTarget(root, target string) (string, error) {
	if strings.TrimSpace(target) == "" {
		return "", fmt.Errorf("target file path is required")
	}

	absolute := target
	if !filepath.IsAbs(absolute) {
		absolute = filepath.Join(root, target)
	}

	info, err := os.Stat(absolute)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("target file does not exist: %s", target)
		}
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("target path is a directory: %s", target)
	}

	relative, err := filepath.Rel(root, absolute)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(relative, "..") {
		return "", fmt.Errorf("target file is outside the repository: %s", target)
	}
	return filepath.ToSlash(relative), nil
}

// AssessTarget builds lock guidance for the validated repository path.
func AssessTarget(root, path string, hasLFS bool) (Assessment, error) {
	assessment := Assessment{Path: filepath.ToSlash(path)}
	fileType, ok := cad.Lookup(filepath.Ext(path))
	if !ok || !fileType.RecommendLocking {
		return assessment, nil
	}

	assessment.LockRecommended = true
	if !hasLFS {
		assessment.ConfigurationWarn = append(assessment.ConfigurationWarn, "Git LFS is not installed or not available on PATH")
		return assessment, nil
	}

	attributesData, err := os.ReadFile(filepath.Join(root, ".gitattributes"))
	if err != nil {
		if os.IsNotExist(err) {
			assessment.ConfigurationWarn = append(assessment.ConfigurationWarn, ".gitattributes is missing the expected Git LFS rule")
			return assessment, nil
		}
		return Assessment{}, err
	}

	if !strings.Contains(string(attributesData), gitx.AttributeLine(fileType.Extension)) {
		assessment.ConfigurationWarn = append(assessment.ConfigurationWarn, fmt.Sprintf("Git LFS is not tracking %s files in .gitattributes", fileType.Extension))
	}

	return assessment, nil
}
