// Package interpolator provides key-path interpolation for merged config maps.
// It resolves cross-references between keys using a ${key.path} syntax,
// allowing config values to reference other values within the same document.
package interpolator

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Interpolate resolves all ${key.path} references in string values within data.
// References are resolved using dot-separated paths into the top-level map.
// Returns an error if a reference cannot be resolved or a cycle is detected.
func Interpolate(data map[string]any) (map[string]any, error) {
	out := shallowCopy(data)
	return interpolateMap(out, out, 0)
}

func interpolateMap(m map[string]any, root map[string]any, depth int) (map[string]any, error) {
	if depth > 32 {
		return nil, fmt.Errorf("interpolation depth limit exceeded (possible cycle)")
	}
	result := make(map[string]any, len(m))
	for k, v := range m {
		resolved, err := interpolateValue(v, root, depth)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

func interpolateValue(v any, root map[string]any, depth int) (any, error) {
	switch val := v.(type) {
	case string:
		return interpolateString(val, root)
	case map[string]any:
		return interpolateMap(val, root, depth+1)
	case []any:
		return interpolateSlice(val, root, depth)
	default:
		return v, nil
	}
}

func interpolateString(s string, root map[string]any) (string, error) {
	var lastErr error
	result := refPattern.ReplaceAllStringFunc(s, func(match string) string {
		path := refPattern.FindStringSubmatch(match)[1]
		val, err := resolvePath(root, strings.Split(path, "."))
		if err != nil {
			lastErr = fmt.Errorf("unresolved reference ${%s}: %w", path, err)
			return match
		}
		return fmt.Sprintf("%v", val)
	})
	if lastErr != nil {
		return "", lastErr
	}
	return result, nil
}

func interpolateSlice(sl []any, root map[string]any, depth int) ([]any, error) {
	out := make([]any, len(sl))
	for i, item := range sl {
		resolved, err := interpolateValue(item, root, depth)
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
		out[i] = resolved
	}
	return out, nil
}

func resolvePath(m map[string]any, parts []string) (any, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty path")
	}
	v, ok := m[parts[0]]
	if !ok {
		return nil, fmt.Errorf("key %q not found", parts[0])
	}
	if len(parts) == 1 {
		return v, nil
	}
	nested, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("key %q is not a map", parts[0])
	}
	return resolvePath(nested, parts[1:])
}

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
