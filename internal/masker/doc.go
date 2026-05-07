// Package masker provides value-level masking for configuration maps.
//
// Unlike the redactor package, which targets keys matching sensitive patterns
// (e.g. "password", "secret"), the masker targets specific string values or
// applies a caller-supplied predicate to decide whether a value should be
// replaced with a placeholder.
//
// Basic usage:
//
//	out := masker.Mask(cfg, masker.Options{
//		Values:      []string{"hunter2", "tok_abc123"},
//		Placeholder: "[MASKED]",
//	})
//
// A MatchFn can be supplied for dynamic decisions:
//
//	out := masker.Mask(cfg, masker.Options{
//		MatchFn: func(key, value string) bool {
//			return strings.HasPrefix(value, "tok_")
//		},
//	})
//
// The original map is never mutated; a deep copy is returned.
package masker
