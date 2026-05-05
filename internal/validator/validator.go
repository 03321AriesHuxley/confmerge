package validator

import (
	"fmt"
	"strings"
)

// ValidationError holds all validation issues found in a config map.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d error(s):\n  - %s",
		len(e.Errors), strings.Join(e.Errors, "\n  - "))
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// Validate checks a merged config map for common structural issues.
// It ensures no keys are empty strings and that nested maps are valid.
func Validate(data map[string]interface{}) error {
	ve := &ValidationError{}
	validateMap(data, "", ve)
	if ve.HasErrors() {
		return ve
	}
	return nil
}

func validateMap(data map[string]interface{}, prefix string, ve *ValidationError) {
	for k, v := range data {
		if k == "" {
			path := prefix
			if path == "" {
				path = "<root>"
			}
			ve.Errors = append(ve.Errors, fmt.Sprintf("empty key found under %q", path))
			continue
		}

		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			validateMap(val, fullKey, ve)
		case map[interface{}]interface{}:
			converted := convertMap(val)
			validateMap(converted, fullKey, ve)
		case nil:
			ve.Errors = append(ve.Errors, fmt.Sprintf("key %q has a null value", fullKey))
		}
	}
}

func convertMap(in map[interface{}]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[fmt.Sprintf("%v", k)] = v
	}
	return out
}
