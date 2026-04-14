package cli

import (
	"fmt"
	"os"

	"github.com/cadops/cadops/internal/doctor"
	"github.com/cadops/cadops/internal/gitx"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Validate repository health",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			return runDoctor(dir)
		},
	}
}

func runDoctor(dir string) error {
	report := doctor.Run(dir, gitx.Runner{})
	for _, result := range report.Results {
		fmt.Printf("[%s] %s: %s\n", result.Level, result.Name, result.Details)
	}
	if report.HasFailures() {
		return fmt.Errorf("doctor found failing checks")
	}
	return nil
}
