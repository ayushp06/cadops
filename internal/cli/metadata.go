package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/metadata"
	"github.com/cadops/cadops/internal/watch"
	"github.com/spf13/cobra"
)

func newMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Generate and inspect CAD file metadata",
	}

	cmd.AddCommand(newMetadataGenerateCmd())
	cmd.AddCommand(newMetadataShowCmd())
	return cmd
}

func newMetadataGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate repository CAD metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runMetadataGenerate(dir)
		},
	}
}

func newMetadataShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <file>",
		Short: "Show stored metadata for a CAD file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runMetadataShow(dir, args[0])
		},
	}
}

func runMetadataGenerate(dir string) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	repoRoot, err := gitx.RepoRoot(runner, dir)
	if err != nil {
		return err
	}

	cfg, err := config.Load(filepath.Join(repoRoot, config.FileName))
	if err != nil {
		return err
	}

	manifest, err := metadata.Generate(repoRoot, cfg.TrackedExtensions)
	if err != nil {
		return err
	}
	if err := metadata.Save(repoRoot, manifest); err != nil {
		return err
	}

	fmt.Println(formatMetadataGenerateReport(len(manifest.Records)))
	return nil
}

func runMetadataShow(dir, target string) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	repoRoot, err := gitx.RepoRoot(runner, dir)
	if err != nil {
		return err
	}

	cfg, err := config.Load(filepath.Join(repoRoot, config.FileName))
	if err != nil {
		return err
	}

	relPath, absPath, err := resolveRepoRelativePath(repoRoot, dir, target)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", target)
		}
		return err
	}

	if !watch.NewFilter(cfg.TrackedExtensions).Match(relPath) {
		return fmt.Errorf("not a configured CAD file: %s", relPath)
	}

	manifest, err := metadata.Load(repoRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("metadata store not found; run `cadops metadata generate`")
		}
		return err
	}

	record, ok := metadata.Lookup(manifest, relPath)
	if !ok {
		return fmt.Errorf("metadata not found for %s; run `cadops metadata generate`", relPath)
	}

	fmt.Print(formatMetadataRecord(record))
	return nil
}

func formatMetadataGenerateReport(count int) string {
	return fmt.Sprintf("Generated metadata for %d CAD files at .cadops/metadata/manifest.json", count)
}

func formatMetadataRecord(record metadata.Record) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Path: %s\n", record.Path))
	builder.WriteString(fmt.Sprintf("Type: %s\n", record.TypeName))
	builder.WriteString(fmt.Sprintf("Extension: %s\n", record.Extension))
	builder.WriteString(fmt.Sprintf("Size: %d bytes\n", record.SizeBytes))
	builder.WriteString(fmt.Sprintf("Modified: %s\n", record.ModifiedTime))
	builder.WriteString(fmt.Sprintf("SHA-256: %s\n", record.SHA256))
	builder.WriteString(fmt.Sprintf("Git LFS Expected: %s\n", yesNo(record.GitLFSExpected)))
	builder.WriteString(fmt.Sprintf("Locking Recommended: %s\n", yesNo(record.LockingRecommended)))
	return builder.String()
}

func resolveRepoRelativePath(repoRoot, currentDir, target string) (string, string, error) {
	absPath := target
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(currentDir, target)
	}

	absPath, err := filepath.Abs(absPath)
	if err != nil {
		return "", "", err
	}

	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		return "", "", err
	}
	relPath = filepath.Clean(relPath)
	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("file is outside the repository: %s", target)
	}

	return filepath.ToSlash(relPath), absPath, nil
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
