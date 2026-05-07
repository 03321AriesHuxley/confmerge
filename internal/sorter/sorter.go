package sorter

import (
	"sort"
	"strings"
)

// Options controls how keys are sorted in the output map.
type Options struct {
	// Recursive sorts nested maps as well.
	Recursive bool
	// CaseInsensitive sorts keys without regard to case.
	CaseInsensitive bool
	// Descending reverses the sort order.
	Descending bool
}

// DefaultOptions returns a sensible default: recursive, case-insensitive,
// ascending sort.
func DefaultOptions() Options {
	return Options{
		Recursive:       true,
		CaseInsensitive: true,
		Descending:      false,
	}
}

// Sort returns a new map whose keys are sorted according to opts. The input
// map is never mutated.
func Sort(input map[string]any, opts Options) map[string]any {
	if input == nil {
		return nil
	}
	out := make(map[string]any, len(input))
	keys := sortedKeys(input, opts)
	for _, k := range keys {
		v := input[k]
		if opts.Recursive {
			if child, ok := toMap(v); ok {
				v = Sort(child, opts)
			}
		}
		out[k] = v
	}
	return out
}

// Keys returns the keys of input in sorted order according to opts.
func Keys(input map[string]any, opts Options) []string {
	return sortedKeys(input, opts)
}

func sortedKeys(m map[string]any, opts Options) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		if opts.CaseInsensitive {
			a = strings.ToLower(a)
			b = strings.ToLower(b)
		}
		less := a < b
		if opts.Descending {
			return !less
		}
		return less
	})
	return keys
}

func toMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	return m, ok
}
