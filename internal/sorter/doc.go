// Package sorter provides utilities for producing deterministically ordered
// representations of config maps.
//
// When config files are serialised to YAML or JSON the key order is
// non-deterministic by default, which makes diffs noisy and hard to review.
// The sorter package solves this by returning a new map whose keys are sorted
// according to caller-supplied Options.
//
// Basic usage:
//
//	out := sorter.Sort(merged, sorter.DefaultOptions())
//
// Sorting is non-destructive: the original map is never modified.
package sorter
