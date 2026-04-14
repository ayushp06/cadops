package cli

import (
	"fmt"
	"os"

	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/history"
	"github.com/spf13/cobra"
)

const defaultHistoryLimit = 10

func newHistoryCmd() *cobra.Command {
	limit := defaultHistoryLimit
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show recent CAD-aware commit history",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runHistory(dir, limit)
		},
	}
	cmd.Flags().IntVarP(&limit, "limit", "n", defaultHistoryLimit, "Number of recent commits to show")
	return cmd
}

func runHistory(dir string, limit int) error {
	if limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}
	if !gitx.IsInstalled("git") {
		return fmt.Errorf("git is not installed or not on PATH")
	}

	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}
	if !gitx.HasCommits(runner, dir) {
		fmt.Println("No commits found")
		return nil
	}

	out, err := gitx.RecentHistory(runner, dir, limit)
	if err != nil {
		return err
	}

	fmt.Print(history.Format(history.ParseLog(out)))
	return nil
}
