package gitx

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const commandTimeout = 10 * time.Second

// Result captures a completed command invocation.
type Result struct {
	Stdout string
	Stderr string
}

// Runner executes git-family commands.
type Runner struct{}

// Run executes a command in the given working directory.
func (Runner) Run(dir, name string, args ...string) (Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := Result{
		Stdout: strings.TrimSpace(stdout.String()),
		Stderr: strings.TrimSpace(stderr.String()),
	}
	if err != nil {
		return result, fmt.Errorf("%s %v: %w: %s", name, args, err, result.Stderr)
	}
	return result, nil
}

// IsInstalled reports whether the command exists on PATH.
func IsInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
