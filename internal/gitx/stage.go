package gitx

// StagePath stages a path, including deletions, in the working tree.
func StagePath(runner Runner, dir, path string) error {
	_, err := runner.Run(dir, "git", "add", "--all", "--", path)
	return err
}
