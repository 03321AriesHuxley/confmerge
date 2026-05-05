// Package transformer provides utilities for transforming merged config maps
// before writing output, such as flattening keys or applying environment
// variable substitution.
package transformer

import (
	"fmt"
	"os"
	"strings"
)

// Options controls which transformations are applied.
type Options struct {
	// EnvSubst replaces ${VAR} and $VAR references with environment variable values.
	EnvSubst bool
	// FlattenSep, when non-empty, flattens nested maps into dot-separated keys.
	FlattenSep string
}

// Apply runs all enabled transformations on data and returns the result.
func Apply(data map[string]any, opts Options) (map[string]any, error) {
	result := data
	if opts.EnvSubst {
		var err error
		result, err = envSubstMap(result)
		if err != nil {
			return nil, fmt.Errorf("env substitution: %w", err)
		}
	}
	if opts.FlattenSep != "" {
		result = flattenMap(result, opts.FlattenSep, "")
	}
	return result, nil
}

// envSubstMap recursively replaces environment variable references in string values.
func envSubstMap(data map[string]any) (map[string]any, error) {
	out := make(map[string]any, len(data))
	for k, v := range data {
		switch val := v.(type) {
		case string:
			out[k] = os.ExpandEnv(val)
		case map[string]any:
			substituted, err := envSubstMap(val)
			if err != nil {
				return nil, err
			}
			out[k] = substituted
		default:
			out[k] = v
		}
	}
	return out, nil
}

// flattenMap recursively flattens nested maps using sep as the key separator.
func flattenMap(data map[string]any, sep, prefix string) map[string]any {
	out := make(map[string]any)
	for k, v := range data {
		fullKey := k
		if prefix != "" {
			fullKey = strings.Join([]string{prefix, k}, sep)
		}
		if nested, ok := v.(map[string]any); ok {
			for fk, fv := range flattenMap(nested, sep, fullKey) {
				out[fk] = fv
			}
		} else {
			out[fullKey] = v
		}
	}
	return out
}
