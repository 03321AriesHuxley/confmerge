// Package renamer provides key renaming functionality for config maps.
// It supports renaming top-level and nested keys using dot-separated paths.
package renamer

import "fmt"

// Rule defines a single rename operation from OldPath to NewPath.
type Rule struct {
	OldPath string
	NewPath string
}

// Rename applies a slice of rename rules to a copy of the provided map.
// Keys are addressed using dot-separated paths (e.g. "database.host").
// Returns an error if a source path does not exist.
func Rename(src map[string]any, rules []Rule) (map[string]any, error) {
	dst := deepCopy(src)
	for _, r := range rules {
		val, ok := getPath(dst, segments(r.OldPath))
		if !ok {
			return nil, fmt.Errorf("renamer: key not found: %s", r.OldPath)
		}
		if err := deletePath(dst, segments(r.OldPath)); err != nil {
			return nil, fmt.Errorf("renamer: delete %s: %w", r.OldPath, err)
		}
		if err := setPath(dst, segments(r.NewPath), val); err != nil {
			return nil, fmt.Errorf("renamer: set %s: %w", r.NewPath, err)
		}
	}
	return dst, nil
}

func segments(path string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '.' {
			parts = append(parts, path[start:i])
			start = i + 1
		}
	}
	parts = append(parts, path[start:])
	return parts
}

func getPath(m map[string]any, segs []string) (any, bool) {
	if len(segs) == 1 {
		v, ok := m[segs[0]]
		return v, ok
	}
	child, ok := m[segs[0]]
	if !ok {
		return nil, false
	}
	nested, ok := child.(map[string]any)
	if !ok {
		return nil, false
	}
	return getPath(nested, segs[1:])
}

func setPath(m map[string]any, segs []string, val any) error {
	if len(segs) == 1 {
		m[segs[0]] = val
		return nil
	}
	child, ok := m[segs[0]]
	if !ok {
		child = map[string]any{}
		m[segs[0]] = child
	}
	nested, ok := child.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map at %q, got %T", segs[0], child)
	}
	return setPath(nested, segs[1:], val)
}

func deletePath(m map[string]any, segs []string) error {
	if len(segs) == 1 {
		delete(m, segs[0])
		return nil
	}
	child, ok := m[segs[0]]
	if !ok {
		return nil
	}
	nested, ok := child.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map at %q, got %T", segs[0], child)
	}
	return deletePath(nested, segs[1:])
}

func deepCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if nested, ok := v.(map[string]any); ok {
			out[k] = deepCopy(nested)
		} else {
			out[k] = v
		}
	}
	return out
}
