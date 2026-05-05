package merger

import "fmt"

// MergeStrategy defines how conflicting keys are resolved.
type MergeStrategy int

const (
	// StrategyOverride replaces the base value with the override value.
	StrategyOverride MergeStrategy = iota
	// StrategyKeep retains the base value when a conflict occurs.
	StrategyKeep
)

// Merge performs a deep merge of src into dst.
// Nested maps are merged recursively; all other values follow the given strategy.
// dst is modified in place and also returned for convenience.
func Merge(dst, src map[string]any, strategy MergeStrategy) (map[string]any, error) {
	for key, srcVal := range src {
		dstVal, exists := dst[key]
		if !exists {
			dst[key] = srcVal
			continue
		}

		srcMap, srcIsMap := toMap(srcVal)
		dstMap, dstIsMap := toMap(dstVal)

		switch {
		case srcIsMap && dstIsMap:
			merged, err := Merge(dstMap, srcMap, strategy)
			if err != nil {
				return nil, fmt.Errorf("merging key %q: %w", key, err)
			}
			dst[key] = merged
		case strategy == StrategyOverride:
			dst[key] = srcVal
		case strategy == StrategyKeep:
			// retain existing dst value — no-op
		default:
			return nil, fmt.Errorf("unknown merge strategy: %d", strategy)
		}
	}
	return dst, nil
}

// toMap attempts to cast v to map[string]any.
func toMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	return m, ok
}
