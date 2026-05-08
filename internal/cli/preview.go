package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/preview"
	"github.com/cadops/cadops/internal/watch"
	"github.com/spf13/cobra"
)

func newPreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Generate and inspect CAD preview records",
	}
	cmd.AddCommand(newPreviewGenerateCmd())
	cmd.AddCommand(newPreviewShowCmd())
	cmd.AddCommand(newPreviewListCmd())
	return cmd
}

func newPreviewGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate repository CAD preview records",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runPreviewGenerate(dir, time.Now())
		},
	}
}

func newPreviewShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <file>",
		Short: "Show stored preview record for a CAD file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runPreviewShow(dir, args[0])
		},
	}
}

func newPreviewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List stored CAD preview records",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runPreviewList(dir)
		},
	}
}

func runPreviewGenerate(dir string, now time.Time) error {
	repoRoot, cfg, err := loadRepoConfig(dir)
	if err != nil {
		return err
	}

	manifest, err := preview.Generate(repoRoot, cfg.TrackedExtensions, now)
	if err != nil {
		return err
	}
	if err := preview.Save(repoRoot, manifest); err != nil {
		return err
	}

	fmt.Print(formatPreviewGenerateReport(manifest))
	return nil
}

func runPreviewShow(dir, target string) error {
	repoRoot, cfg, err := loadRepoConfig(dir)
	if err != nil {
		return err
	}

	relPath, _, err := resolveRepoRelativePath(repoRoot, dir, target)
	if err != nil {
		return err
	}
	if !watch.NewFilter(cfg.TrackedExtensions).Match(relPath) {
		return fmt.Errorf("not a configured CAD file: %s", relPath)
	}

	manifest, err := preview.Load(repoRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("preview store not found; run `cadops preview generate`")
		}
		return err
	}

	record, ok := preview.Lookup(manifest, relPath)
	if !ok {
		return fmt.Errorf("preview not found for %s; run `cadops preview generate`", relPath)
	}

	fmt.Print(formatPreviewRecord(repoRoot, record))
	return nil
}

func runPreviewList(dir string) error {
	repoRoot, _, err := loadRepoConfig(dir)
	if err != nil {
		return err
	}

	manifest, err := preview.Load(repoRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("preview store not found; run `cadops preview generate`")
		}
		return err
	}

	fmt.Print(formatPreviewList(repoRoot, manifest))
	return nil
}

func loadRepoConfig(dir string) (string, config.Config, error) {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return "", config.Config{}, fmt.Errorf("not a git repository")
	}

	repoRoot, err := gitx.RepoRoot(runner, dir)
	if err != nil {
		return "", config.Config{}, err
	}

	cfg, err := config.Load(filepath.Join(repoRoot, config.FileName))
	if err != nil {
		return "", config.Config{}, err
	}
	return repoRoot, cfg, nil
}

func formatPreviewGenerateReport(manifest preview.Manifest) string {
	counts := previewStatusCounts("", manifest)
	return fmt.Sprintf(
		"Generated preview records for %d CAD files at .cadops/previews/manifest.json\nGenerated: %d\nUnavailable: %d\nUnsupported: %d\n",
		len(manifest.Records),
		counts[preview.StatusGenerated],
		counts[preview.StatusUnavailable],
		counts[preview.StatusUnsupported],
	)
}

func formatPreviewList(root string, manifest preview.Manifest) string {
	if len(manifest.Records) == 0 {
		return "No preview records found\n"
	}

	var builder strings.Builder
	builder.WriteString("Preview Records\n")
	for _, record := range manifest.Records {
		builder.WriteString(fmt.Sprintf("- %s | %s | %s", record.SourcePath, record.CADType, preview.EffectiveStatus(root, record)))
		if record.ArtifactPath != "" {
			builder.WriteString(fmt.Sprintf(" | artifact: %s", record.ArtifactPath))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func formatPreviewRecord(root string, record preview.Record) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Source: %s\n", record.SourcePath))
	builder.WriteString(fmt.Sprintf("CAD Type: %s\n", record.CADType))
	builder.WriteString(fmt.Sprintf("Source SHA-256: %s\n", record.SourceHash))
	builder.WriteString(fmt.Sprintf("Status: %s\n", preview.EffectiveStatus(root, record)))
	if record.ArtifactPath != "" {
		builder.WriteString(fmt.Sprintf("Artifact: %s\n", record.ArtifactPath))
	}
	builder.WriteString(fmt.Sprintf("Generated: %s\n", record.GeneratedAt))
	if record.Message != "" {
		builder.WriteString(fmt.Sprintf("Message: %s\n", record.Message))
	}
	return builder.String()
}

func previewStatusCounts(root string, manifest preview.Manifest) map[preview.Status]int {
	counts := map[preview.Status]int{
		preview.StatusGenerated:   0,
		preview.StatusUnavailable: 0,
		preview.StatusUnsupported: 0,
		preview.StatusStale:       0,
	}
	for _, record := range manifest.Records {
		status := record.Status
		if root != "" {
			status = preview.EffectiveStatus(root, record)
		}
		counts[status]++
	}
	return counts
}
