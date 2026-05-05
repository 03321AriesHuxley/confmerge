// Package differ provides utilities for computing the difference between
// two configuration maps (base vs override).
//
// It produces a list of Entry values, each describing whether a key was
// added, removed, modified, or unchanged when comparing the base config
// to the override config. Nested maps are traversed recursively and
// keys are reported using dot-notation (e.g. "database.host").
//
// Example usage:
//
//	base := map[string]interface{}{"host": "localhost", "port": 5432}
//	override := map[string]interface{}{"host": "remotehost", "port": 5432}
//	entries := differ.Diff(base, override)
//	for _, e := range entries {
//		fmt.Printf("%s: %s\n", e.Key, e.Type)
//	}
package differ
