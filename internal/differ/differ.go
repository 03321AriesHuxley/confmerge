package differ

import (
	"fmt"
	"sort"
)

// DiffType represents the type of difference between two config values.
type DiffType string

const (
	Added    DiffType = "added"
	Removed  DiffType = "removed"
	Modified DiffType = "modified"
	Unchanged DiffType = "unchanged"
)

// Entry represents a single difference between two config maps.
type Entry struct {
	Key      string
	Type     DiffType
	OldValue interface{}
	NewValue interface{}
}

// Diff computes the difference between two config maps (base vs override).
// It returns a slice of Entry describing each change.
func Diff(base, override map[string]interface{}) []Entry {
	return diffMaps("", base, override)
}

func diffMaps(prefix string, base, override map[string]interface{}) []Entry {
	var entries []Entry
	seen := make(map[string]bool)

	keys := sortedKeys(base)
	for _, k := range keys {
		fullKey := joinKey(prefix, k)
		seen[k] = true
		baseVal := base[k]
		overrideVal, exists := override[k]
		if !exists {
			entries = append(entries, Entry{Key: fullKey, Type: Removed, OldValue: baseVal})
			continue
		}
		baseMap, baseIsMap := toMap(baseVal)
		overrideMap, overrideIsMap := toMap(overrideVal)
		if baseIsMap && overrideIsMap {
			entries = append(entries, diffMaps(fullKey, baseMap, overrideMap)...)
		} else if fmt.Sprintf("%v", baseVal) != fmt.Sprintf("%v", overrideVal) {
			entries = append(entries, Entry{Key: fullKey, Type: Modified, OldValue: baseVal, NewValue: overrideVal})
		} else {
			entries = append(entries, Entry{Key: fullKey, Type: Unchanged, OldValue: baseVal, NewValue: overrideVal})
		}
	}

	for _, k := range sortedKeys(override) {
		if seen[k] {
			continue
		}
		fullKey := joinKey(prefix, k)
		entries = append(entries, Entry{Key: fullKey, Type: Added, NewValue: override[k]})
	}
	return entries
}

func joinKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func toMap(v interface{}) (map[string]interface{}, bool) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, true
	}
	return nil, false
}
