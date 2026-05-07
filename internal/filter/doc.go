// Package filter provides key-based include/exclude filtering for
// configuration maps.
//
// Rules are expressed as dot-separated key paths (e.g. "database.host").
// When an Include list is supplied only matching paths are retained;
// when an Exclude list is supplied those paths are dropped from the
// output. Both lists may be combined: exclusions are applied first.
//
// Nested maps are traversed recursively and the original input is
// never mutated.
package filter
