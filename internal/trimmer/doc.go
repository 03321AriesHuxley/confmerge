// Package trimmer provides utilities for cleaning up map keys and string
// values within a configuration map.
//
// It supports trimming leading/trailing whitespace from keys and string
// values, as well as removing configurable prefixes and suffixes from
// string values. All operations produce a deep copy and never mutate
// the original input.
//
// Basic usage:
//
//	out := trimmer.Trim(cfg, trimmer.DefaultOptions())
//
// Custom options:
//
//	out := trimmer.Trim(cfg, trimmer.Options{
//		TrimKeys:         true,
//		TrimStringValues: true,
//		TrimPrefix:       "env:",
//	})
package trimmer
