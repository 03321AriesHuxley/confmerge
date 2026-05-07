package filter

import "strings"

// Options controls which keys are included or excluded during filtering.
type Options struct {
	// Include is a list of dot-separated key paths to include (whitelist).
	// If non-empty, only matching keys are retained.
	Include []string

	// Exclude is a list of dot-separated key paths to exclude (blacklist).
	Exclude []string
}

// Filter returns a new map containing only the keys that satisfy the
// include/exclude rules defined in opts. Nested maps are traversed
// recursively. The original map is never mutated.
func Filter(data map[string]any, opts Options) map[string]any {
	includeSet := toSet(opts.Include)
	excludeSet := toSet(opts.Exclude)
	return filterMap(data, includeSet, excludeSet, "")
}

func filterMap(data map[string]any, include, exclude map[string]struct{}, prefix string) map[string]any {
	out := make(map[string]any, len(data))
	for k, v := range data {
		path := join(prefix, k)

		if _, ok := exclude[path]; ok {
			continue
		}

		if len(include) > 0 {
			if !isIncluded(path, include) {
				continue
			}
		}

		if nested, ok := v.(map[string]any); ok {
			v = filterMap(nested, include, exclude, path)
		}

		out[k] = v
	}
	return out
}

// isIncluded returns true when path matches or is a prefix of any include entry.
func isIncluded(path string, include map[string]struct{}) bool {
	for entry := range include {
		if entry == path || strings.HasPrefix(entry, path+".") || strings.HasPrefix(path, entry+".") {
			return true
		}
	}
	return false
}

func join(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
