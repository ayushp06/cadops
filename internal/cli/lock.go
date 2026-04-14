package cli

import (
	"fmt"
	"os"

	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/locking"
	"github.com/spf13/cobra"
)

func newLockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lock <file>",
		Short: "Lock a CAD file with Git LFS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runLock(dir, args[0])
		},
	}
}

func newUnlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlock <file>",
		Short: "Unlock a CAD file with Git LFS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runUnlock(dir, args[0])
		},
	}
}

func runLock(dir, target string) error {
	return runLockAction(dir, target, "Locked", gitx.LockPath)
}

func runUnlock(dir, target string) error {
	return runLockAction(dir, target, "Unlocked", gitx.UnlockPath)
}

func runLockAction(dir, target, verb string, action func(gitx.Runner, string, string) error) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	path, err := locking.ResolveTarget(dir, target)
	if err != nil {
		return err
	}

	assessment, err := locking.AssessTarget(dir, path, gitx.HasLFS(runner, dir))
	if err != nil {
		return err
	}
	for _, warning := range assessment.ConfigurationWarn {
		fmt.Printf("Warning: %s\n", warning)
	}

	if err := action(runner, dir, path); err != nil {
		return err
	}

	fmt.Printf("%s %s\n", verb, path)
	return nil
}
