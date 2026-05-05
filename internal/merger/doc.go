// Package merger provides deep-merge functionality for layered configuration maps.
//
// It supports two merge strategies:
//
//   - StrategyOverride: values from the source (higher-precedence layer) overwrite
//     conflicting scalar values in the destination.
//
//   - StrategyKeep: the destination's existing scalar values are preserved when a
//     conflict is detected.
//
// In both strategies, nested maps are always merged recursively so that only the
// conflicting leaf values are subject to the chosen strategy.
//
// Typical usage:
//
//	base := map[string]any{"log": map[string]any{"level": "info"}}
//	override := map[string]any{"log": map[string]any{"level": "debug"}}
//
//	result, err := merger.Merge(base, override, merger.StrategyOverride)
//	// result["log"].(map[string]any)["level"] == "debug"
package merger
