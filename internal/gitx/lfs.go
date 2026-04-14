package gitx

import (
	"strings"
)

// EnsureLFSInstalled runs `git lfs install` when the environment supports it.
func EnsureLFSInstalled(runner Runner, dir string) error {
	_, err := runner.Run(dir, "git", "lfs", "install")
	return err
}

// ListTrackedPatterns returns the patterns from `git lfs track --list`.
func ListTrackedPatterns(runner Runner, dir string) ([]string, error) {
	result, err := runner.Run(dir, "git", "lfs", "track", "--list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(result.Stdout, "\n")
	patterns := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		patterns = append(patterns, strings.Trim(parts[0], "\""))
	}
	return patterns, nil
}
