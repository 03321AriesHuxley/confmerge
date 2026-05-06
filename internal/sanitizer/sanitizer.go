// Package sanitizer provides utilities for normalizing and sanitizing
// merged configuration maps before further processing or output.
//
// It trims whitespace from string values, removes nil/null entries,
// and optionally coerces numeric strings to their native types.
package sanitizer

import (
	"fmt"
	"strconv"
	"strings"
)

// Options controls which sanitization passes are applied.
type Options struct {
	// TrimStrings removes leading/trailing whitespace from all string values.
	TrimStrings bool
	// DropNulls removes keys whose value is nil.
	DropNulls bool
	// CoerceNumbers attempts to parse string values as int64 or float64.
	CoerceNumbers bool
}

// DefaultOptions returns a sensible default sanitizer configuration.
func DefaultOptions() Options {
	return Options{
		TrimStrings:   true,
		DropNulls:     true,
		CoerceNumbers: false,
	}
}

// Sanitize applies the given Options to a configuration map and returns
// a new, sanitized copy without mutating the original.
func Sanitize(input map[string]any, opts Options) map[string]any {
	return sanitizeMap(input, opts)
}

func sanitizeMap(m map[string]any, opts Options) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		sanitized := sanitizeValue(v, opts)
		if opts.DropNulls && sanitized == nil {
			continue
		}
		out[k] = sanitized
	}
	return out
}

func sanitizeValue(v any, opts Options) any {
	switch val := v.(type) {
	case map[string]any:
		return sanitizeMap(val, opts)
	case string:
		result := val
		if opts.TrimStrings {
			result = strings.TrimSpace(result)
		}
		if opts.CoerceNumbers {
			if coerced, err := coerce(result); err == nil {
				return coerced
			}
		}
		return result
	case []any:
		return sanitizeSlice(val, opts)
	case nil:
		return nil
	default:
		return v
	}
}

func sanitizeSlice(s []any, opts Options) []any {
	out := make([]any, 0, len(s))
	for _, item := range s {
		sanitized := sanitizeValue(item, opts)
		if opts.DropNulls && sanitized == nil {
			continue
		}
		out = append(out, sanitized)
	}
	return out
}

func coerce(s string) (any, error) {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("not a number: %q", s)
}
