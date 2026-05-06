// Package patcher applies patch operations (set, delete, merge) to a
// config map using a simple dot-notation path syntax.
package patcher

import (
	"fmt"
	"strings"
)

// Op is the type of patch operation.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpMerge  Op = "merge"
)

// Patch describes a single patch operation.
type Patch struct {
	Op    Op
	Path  string
	Value interface{}
}

// Apply applies a slice of patches to the given config map.
// The map is mutated in place and also returned for convenience.
func Apply(cfg map[string]interface{}, patches []Patch) (map[string]interface{}, error) {
	for _, p := range patches {
		var err error
		switch p.Op {
		case OpSet:
			err = setPath(cfg, p.Path, p.Value)
		case OpDelete:
			err = deletePath(cfg, p.Path)
		case OpMerge:
			src, ok := p.Value.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("patch merge at %q: value must be a map", p.Path)
			}
			err = mergePath(cfg, p.Path, src)
		default:
			return nil, fmt.Errorf("unknown patch op: %q", p.Op)
		}
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func segments(path string) []string {
	return strings.Split(strings.Trim(path, "."), ".")
}

func setPath(cfg map[string]interface{}, path string, value interface{}) error {
	parts := segments(path)
	m := cfg
	for i, key := range parts[:len(parts)-1] {
		v, ok := m[key]
		if !ok {
			next := make(map[string]interface{})
			m[key] = next
			m = next
			continue
		}
		next, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("set %q: key %q at depth %d is not a map", path, key, i)
		}
		m = next
	}
	m[parts[len(parts)-1]] = value
	return nil
}

func deletePath(cfg map[string]interface{}, path string) error {
	parts := segments(path)
	m := cfg
	for i, key := range parts[:len(parts)-1] {
		v, ok := m[key]
		if !ok {
			return nil // already absent
		}
		next, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("delete %q: key %q at depth %d is not a map", path, key, i)
		}
		m = next
	}
	delete(m, parts[len(parts)-1])
	return nil
}

func mergePath(cfg map[string]interface{}, path string, src map[string]interface{}) error {
	if path == "" || path == "." {
		deepMerge(cfg, src)
		return nil
	}
	var target interface{}
	parts := segments(path)
	m := cfg
	for i, key := range parts[:len(parts)-1] {
		v, ok := m[key]
		if !ok {
			next := make(map[string]interface{})
			m[key] = next
			m = next
			continue
		}
		next, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("merge %q: key %q at depth %d is not a map", path, key, i)
		}
		m = next
	}
	last := parts[len(parts)-1]
	target = m[last]
	if target == nil {
		m[last] = src
		return nil
	}
	dst, ok := target.(map[string]interface{})
	if !ok {
		return fmt.Errorf("merge %q: target is not a map", path)
	}
	deepMerge(dst, src)
	return nil
}

func deepMerge(dst, src map[string]interface{}) {
	for k, v := range src {
		if dv, ok := dst[k]; ok {
			if dm, ok := dv.(map[string]interface{}); ok {
				if sm, ok := v.(map[string]interface{}); ok {
					deepMerge(dm, sm)
					continue
				}
			}
		}
		dst[k] = v
	}
}
