package gitx

// CommitPaths creates a commit from the specified paths only.
func CommitPaths(runner Runner, dir, message string, paths []string) error {
	args := []string{"commit", "-m", message, "--only", "--"}
	args = append(args, paths...)
	_, err := runner.Run(dir, "git", args...)
	return err
}
