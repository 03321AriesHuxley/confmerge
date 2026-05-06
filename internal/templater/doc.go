// Package templater implements Go-template rendering for confmerge config values.
//
// After all config layers have been merged, Render can be called to expand any
// string values that contain Go template expressions. Top-level keys of the
// merged map are exposed as template data, enabling cross-key references:
//
//	base_url: "https://api.example.com"
//	health:   "{{ index . \"base_url\" }}/health"
//
// Nested maps are processed recursively. Non-string values (integers, booleans,
// slices, etc.) are passed through without modification.
//
// Rendering uses the "missingkey=error" option so that references to undefined
// keys surface as clear errors rather than silent empty strings.
package templater
