package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a CAD-safe Git repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runInit(dir)
		},
	}
}

func runInit(dir string) error {
	runner := gitx.Runner{}

	if !gitx.IsInstalled("git") {
		return fmt.Errorf("git is not installed or not on PATH")
	}

	if _, err := runner.Run(dir, "git", "lfs", "version"); err != nil {
		return fmt.Errorf("git lfs is not installed or not on PATH")
	}

	repoInitialized := gitx.IsRepo(runner, dir)
	if !repoInitialized {
		if err := gitx.InitRepo(runner, dir); err != nil {
			return err
		}
	}

	if err := gitx.EnsureLFSInstalled(runner, dir); err != nil {
		return err
	}

	cfg := config.Default()
	configPath := filepath.Join(dir, config.FileName)
	if err := config.Save(configPath, cfg); err != nil {
		return err
	}

	if err := ensureAttributes(dir, cfg.TrackedExtensions); err != nil {
		return err
	}
	if err := ensureGitIgnore(dir); err != nil {
		return err
	}

	fmt.Println("CadOps repository initialized")
	if repoInitialized {
		fmt.Println("- Git repository: existing")
	} else {
		fmt.Println("- Git repository: created")
	}
	fmt.Println("- Git LFS: verified and installed")
	fmt.Printf("- Config: wrote %s\n", config.FileName)
	fmt.Println("- Attributes: updated .gitattributes")
	fmt.Println("- Ignore rules: updated .gitignore")
	return nil
}
