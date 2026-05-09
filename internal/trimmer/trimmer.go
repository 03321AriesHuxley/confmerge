package trimmer

import (
	"strings"
)

// Options controls which trimming operations are applied.
type Options struct {
	// TrimKeys removes leading/trailing whitespace from map keys.
	TrimKeys bool
	// TrimStringValues removes leading/trailing whitespace from string values.
	TrimStringValues bool
	// TrimPrefix removes a prefix from all string values if present.
	TrimPrefix string
	// TrimSuffix removes a suffix from all string values if present.
	TrimSuffix string
}

// DefaultOptions returns Options with all basic trimming enabled.
func DefaultOptions() Options {
	return Options{
		TrimKeys:         true,
		TrimStringValues: true,
	}
}

// Trim applies the given Options to a deep copy of src and returns the result.
func Trim(src map[string]any, opts Options) map[string]any {
	return trimMap(src, opts)
}

func trimMap(m map[string]any, opts Options) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		key := k
		if opts.TrimKeys {
			key = strings.TrimSpace(k)
		}
		out[key] = trimValue(v, opts)
	}
	return out
}

func trimValue(v any, opts Options) any {
	switch val := v.(type) {
	case string:
		return trimString(val, opts)
	case map[string]any:
		return trimMap(val, opts)
	case []any:
		return trimSlice(val, opts)
	default:
		return v
	}
}

func trimString(s string, opts Options) string {
	if opts.TrimStringValues {
		s = strings.TrimSpace(s)
	}
	if opts.TrimPrefix != "" {
		s = strings.TrimPrefix(s, opts.TrimPrefix)
	}
	if opts.TrimSuffix != "" {
		s = strings.TrimSuffix(s, opts.TrimSuffix)
	}
	return s
}

func trimSlice(sl []any, opts Options) []any {
	out := make([]any, len(sl))
	for i, item := range sl {
		out[i] = trimValue(item, opts)
	}
	return out
}
