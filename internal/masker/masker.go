// Package masker provides functionality for masking specific values
// in a configuration map, replacing them with a fixed placeholder string.
// Unlike the redactor (which targets sensitive keys), the masker targets
// specific values or applies a custom match function.
package masker

import "strings"

const DefaultPlaceholder = "***"

// Options controls masking behaviour.
type Options struct {
	// Placeholder replaces masked values. Defaults to DefaultPlaceholder.
	Placeholder string
	// Values is a list of exact string values to mask.
	Values []string
	// MatchFn, if non-nil, is called for each string value; returning true
	// causes the value to be masked. Applied after Values matching.
	MatchFn func(key, value string) bool
}

// Mask walks m recursively and returns a new map with matching values
// replaced by opts.Placeholder. The original map is never mutated.
func Mask(m map[string]any, opts Options) map[string]any {
	if opts.Placeholder == "" {
		opts.Placeholder = DefaultPlaceholder
	}
	valSet := make(map[string]struct{}, len(opts.Values))
	for _, v := range opts.Values {
		valSet[v] = struct{}{}
	}
	return maskMap(m, opts, valSet)
}

func maskMap(m map[string]any, opts Options, valSet map[string]struct{}) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = maskValue(k, v, opts, valSet)
	}
	return out
}

func maskValue(key string, v any, opts Options, valSet map[string]struct{}) any {
	switch val := v.(type) {
	case map[string]any:
		return maskMap(val, opts, valSet)
	case string:
		if shouldMask(key, val, opts, valSet) {
			return opts.Placeholder
		}
		return val
	case []any:
		result := make([]any, len(val))
		for i, elem := range val {
			result[i] = maskValue(key, elem, opts, valSet)
		}
		return result
	default:
		return v
	}
}

func shouldMask(key, value string, opts Options, valSet map[string]struct{}) bool {
	if _, ok := valSet[value]; ok {
		return true
	}
	if opts.MatchFn != nil && opts.MatchFn(key, value) {
		return true
	}
	_ = strings.ToLower // imported for potential future use
	return false
}
