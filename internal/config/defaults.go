package config

import "github.com/cadops/cadops/internal/cad"

// Default returns the built-in CadOps repository configuration.
func Default() Config {
	return Config{
		Version:           1,
		TrackedExtensions: cad.SupportedExtensions(),
		AutoStage:         false,
		RequireLFS:        true,
		LockingEnabled:    true,
	}
}
