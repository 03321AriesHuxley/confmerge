// Package linter provides config map linting to detect common
// misconfigurations such as duplicate keys (after case folding),
// deeply nested structures, and suspiciously large scalar values.
package linter

import (
	"fmt"
	"strings"
)

// Issue represents a single linting finding.
type Issue struct {
	Level   string // "warn" or "error"
	Path    string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(i.Level), i.Path, i.Message)
}

// Options controls linter behaviour.
type Options struct {
	MaxDepth       int // default 10; 0 means unlimited
	MaxStringBytes int // default 4096; 0 means unlimited
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxDepth:       10,
		MaxStringBytes: 4096,
	}
}

// Lint walks the merged config map and returns a slice of Issues.
// An empty slice means no issues were found.
func Lint(data map[string]any, opts Options) []Issue {
	var issues []Issue
	lintMap(data, "", 1, opts, &issues)
	return issues
}

func lintMap(m map[string]any, prefix string, depth int, opts Options, issues *[]Issue) {
	if opts.MaxDepth > 0 && depth > opts.MaxDepth {
		*issues = append(*issues, Issue{
			Level:   "warn",
			Path:    prefix,
			Message: fmt.Sprintf("nesting depth %d exceeds maximum %d", depth, opts.MaxDepth),
		})
		return
	}

	// Detect case-folded duplicate keys.
	seen := make(map[string]string, len(m))
	for k := range m {
		lower := strings.ToLower(k)
		if orig, exists := seen[lower]; exists {
			*issues = append(*issues, Issue{
				Level:   "error",
				Path:    join(prefix, k),
				Message: fmt.Sprintf("duplicate key (case-insensitive) conflicts with %q", orig),
			})
		} else {
			seen[lower] = k
		}
	}

	for k, v := range m {
		path := join(prefix, k)
		switch val := v.(type) {
		case map[string]any:
			lintMap(val, path, depth+1, opts, issues)
		case string:
			if opts.MaxStringBytes > 0 && len(val) > opts.MaxStringBytes {
				*issues = append(*issues, Issue{
					Level:   "warn",
					Path:    path,
					Message: fmt.Sprintf("string value length %d exceeds maximum %d bytes", len(val), opts.MaxStringBytes),
				})
			}
		}
	}
}

func join(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}
