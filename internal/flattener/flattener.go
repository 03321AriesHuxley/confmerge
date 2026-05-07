// Package flattener provides utilities for flattening and unflattening
// nested configuration maps into dot-separated key-value pairs.
package flattener

import (
	"fmt"
	"strings"
)

// Options controls the behaviour of the flattener.
type Options struct {
	// Separator is the string used to join key segments. Defaults to ".".
	Separator string
	// MaxDepth limits how deep flattening recurses. 0 means unlimited.
	MaxDepth int
}

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		Separator: ".",
		MaxDepth:  0,
	}
}

// Flatten converts a nested map[string]any into a flat map whose keys are
// dot-separated paths (or separated by opts.Separator).
func Flatten(src map[string]any, opts Options) map[string]any {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	out := make(map[string]any)
	flattenMap(src, "", 0, opts, out)
	return out
}

func flattenMap(m map[string]any, prefix string, depth int, opts Options, out map[string]any) {
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + opts.Separator + k
		}
		if nested, ok := v.(map[string]any); ok && (opts.MaxDepth == 0 || depth < opts.MaxDepth) {
			flattenMap(nested, fullKey, depth+1, opts, out)
		} else {
			out[fullKey] = v
		}
	}
}

// Unflatten converts a flat map with dot-separated keys back into a nested
// map[string]any. It returns an error if two keys produce a conflicting path.
func Unflatten(src map[string]any, opts Options) (map[string]any, error) {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	out := make(map[string]any)
	for k, v := range src {
		parts := strings.Split(k, opts.Separator)
		if err := setNested(out, parts, v); err != nil {
			return nil, fmt.Errorf("flattener: conflict at key %q: %w", k, err)
		}
	}
	return out, nil
}

func setNested(m map[string]any, parts []string, v any) error {
	if len(parts) == 1 {
		m[parts[0]] = v
		return nil
	}
	child, exists := m[parts[0]]
	if !exists {
		next := make(map[string]any)
		m[parts[0]] = next
		return setNested(next, parts[1:], v)
	}
	next, ok := child.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map at segment %q, got %T", parts[0], child)
	}
	return setNested(next, parts[1:], v)
}
