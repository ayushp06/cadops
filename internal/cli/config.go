package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cadops/cadops/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect CadOps repository configuration",
	}

	cmd.AddCommand(newConfigShowCmd())
	cmd.AddCommand(newConfigGetCmd())
	return cmd
}

func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show the current .cadops.yaml values",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runConfigShow(dir)
		},
	}
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a single CadOps configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runConfigGet(dir, args[0])
		},
	}
}

func runConfigShow(dir string) error {
	cfg, err := config.Load(filepath.Join(dir, config.FileName))
	if err != nil {
		return err
	}

	for _, key := range []string{"version", "tracked_extensions", "auto_stage", "require_lfs", "locking_enabled"} {
		value, err := config.Lookup(cfg, key)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", key, config.FormatValue(value))
	}

	return nil
}

func runConfigGet(dir, key string) error {
	cfg, err := config.Load(filepath.Join(dir, config.FileName))
	if err != nil {
		return err
	}

	value, err := config.Lookup(cfg, key)
	if err != nil {
		return err
	}

	fmt.Println(config.FormatValue(value))
	return nil
}
