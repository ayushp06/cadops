package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/cadops/cadops/internal/config"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/cadops/cadops/internal/watch"
	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Watch CAD files and react to repository changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runWatch(dir)
		},
	}
}

func runWatch(dir string) error {
	runner := gitx.Runner{}
	if !gitx.IsRepo(runner, dir) {
		return fmt.Errorf("not a git repository")
	}

	cfg, err := config.Load(filepath.Join(dir, config.FileName))
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Printf("Watching %s for CAD changes\n", dir)
	if cfg.AutoStage {
		fmt.Println("Auto-stage enabled")
	} else {
		fmt.Println("Auto-stage disabled")
	}

	watcher := watch.New(dir, cfg.TrackedExtensions)
	return watcher.Run(ctx, func(event watch.Event) {
		status := fmt.Sprintf("%s %s", event.Kind, event.Path)
		if cfg.AutoStage {
			if err := gitx.StagePath(runner, dir, event.Path); err != nil {
				fmt.Printf("stage error %s: %v\n", event.Path, err)
				return
			}
			status += " staged"
		}
		fmt.Println(status)
	})
}
