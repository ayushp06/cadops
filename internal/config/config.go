package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const FileName = ".cadops.yaml"

// Config is the persisted CadOps repository configuration.
type Config struct {
	Version           int
	TrackedExtensions []string
	AutoStage         bool
	RequireLFS        bool
	LockingEnabled    bool
}

// Load reads the repository configuration from disk.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return Parse(data)
}

// Save writes the configuration using the simple project YAML format.
func Save(path string, cfg Config) error {
	return os.WriteFile(path, []byte(Marshal(cfg)), 0o644)
}

// Parse converts the constrained YAML configuration into Config.
func Parse(data []byte) (Config, error) {
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	cfg := Config{}
	var inExtensions bool

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "- ") {
			if !inExtensions {
				return Config{}, fmt.Errorf("unexpected list item: %q", line)
			}
			cfg.TrackedExtensions = append(cfg.TrackedExtensions, strings.TrimSpace(strings.TrimPrefix(line, "- ")))
			continue
		}

		inExtensions = false
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("invalid config line: %q", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "version":
			number, err := strconv.Atoi(value)
			if err != nil {
				return Config{}, fmt.Errorf("parse version: %w", err)
			}
			cfg.Version = number
		case "tracked_extensions":
			inExtensions = true
		case "auto_stage":
			flag, err := parseBool(value)
			if err != nil {
				return Config{}, fmt.Errorf("parse auto_stage: %w", err)
			}
			cfg.AutoStage = flag
		case "require_lfs":
			flag, err := parseBool(value)
			if err != nil {
				return Config{}, fmt.Errorf("parse require_lfs: %w", err)
			}
			cfg.RequireLFS = flag
		case "locking_enabled":
			flag, err := parseBool(value)
			if err != nil {
				return Config{}, fmt.Errorf("parse locking_enabled: %w", err)
			}
			cfg.LockingEnabled = flag
		default:
			return Config{}, fmt.Errorf("unknown config key: %s", key)
		}
	}

	if cfg.Version == 0 {
		return Config{}, errors.New("missing version")
	}

	return cfg, nil
}

// Marshal renders the constrained YAML configuration in stable order.
func Marshal(cfg Config) string {
	var builder strings.Builder
	builder.WriteString("version: ")
	builder.WriteString(strconv.Itoa(cfg.Version))
	builder.WriteString("\n")
	builder.WriteString("tracked_extensions:\n")
	for _, extension := range cfg.TrackedExtensions {
		builder.WriteString("  - ")
		builder.WriteString(extension)
		builder.WriteString("\n")
	}
	builder.WriteString("auto_stage: ")
	builder.WriteString(strconv.FormatBool(cfg.AutoStage))
	builder.WriteString("\n")
	builder.WriteString("require_lfs: ")
	builder.WriteString(strconv.FormatBool(cfg.RequireLFS))
	builder.WriteString("\n")
	builder.WriteString("locking_enabled: ")
	builder.WriteString(strconv.FormatBool(cfg.LockingEnabled))
	builder.WriteString("\n")
	return builder.String()
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean %q", value)
	}
}
