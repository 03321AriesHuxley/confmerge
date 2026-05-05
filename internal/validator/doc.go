// Package validator provides structural validation for merged configuration maps.
//
// After merging multiple config layers, it is useful to verify the resulting
// map is well-formed before writing it to an output file. This package checks
// for common issues such as null values and empty string keys that may indicate
// a misconfigured input file.
//
// Basic usage:
//
//	if err := validator.Validate(merged); err != nil {
//		log.Fatalf("invalid config: %v", err)
//	}
//
// Validation errors are returned as a *ValidationError which contains a list
// of all individual issues found, allowing callers to report them all at once.
package validator
