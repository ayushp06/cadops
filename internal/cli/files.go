package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/gitx"
)

const gitignoreTemplate = `# CAD and workstation artifacts
*.bak
*.tmp
*.swp
*.lock

# CadOps
.cadops-cache/
`

func ensureAttributes(dir string, extensions []string) error {
	path := filepath.Join(dir, ".gitattributes")
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	merged := gitx.MergeAttributes(string(data), extensions)
	return os.WriteFile(path, []byte(merged), 0o644)
}

func ensureGitIgnore(dir string) error {
	path := filepath.Join(dir, ".gitignore")
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	existing := strings.ReplaceAll(string(data), "\r\n", "\n")
	additions := strings.ReplaceAll(gitignoreTemplate, "\r\n", "\n")
	if existing == "" {
		return os.WriteFile(path, []byte(additions), 0o644)
	}

	lines := strings.Split(existing, "\n")
	seen := make(map[string]bool, len(lines))
	for _, line := range lines {
		seen[strings.TrimSpace(line)] = true
	}

	output := strings.TrimRight(existing, "\n")
	for _, line := range strings.Split(additions, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		output += "\n" + line
	}
	output = strings.TrimRight(output, "\n") + "\n"
	return os.WriteFile(path, []byte(output), 0o644)
}
