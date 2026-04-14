package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Value is a single config lookup result.
type Value struct {
	Key    string
	Scalar string
	List   []string
}

// IsList reports whether the value is list-shaped.
func (v Value) IsList() bool {
	return v.List != nil
}

// Lookup returns a supported configuration value by key.
func Lookup(cfg Config, key string) (Value, error) {
	switch key {
	case "version":
		return Value{Key: key, Scalar: strconv.Itoa(cfg.Version)}, nil
	case "tracked_extensions":
		out := make([]string, len(cfg.TrackedExtensions))
		copy(out, cfg.TrackedExtensions)
		return Value{Key: key, List: out}, nil
	case "auto_stage":
		return Value{Key: key, Scalar: strconv.FormatBool(cfg.AutoStage)}, nil
	case "require_lfs":
		return Value{Key: key, Scalar: strconv.FormatBool(cfg.RequireLFS)}, nil
	case "locking_enabled":
		return Value{Key: key, Scalar: strconv.FormatBool(cfg.LockingEnabled)}, nil
	default:
		return Value{}, fmt.Errorf("unknown config key %q", key)
	}
}

// FormatValue renders a lookup value in a stable CLI-friendly format.
func FormatValue(value Value) string {
	if value.IsList() {
		return strings.Join(value.List, ", ")
	}
	return value.Scalar
}
