package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FieldType represents the expected type of a config field.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeFloat   FieldType = "float"
	TypeMap     FieldType = "map"
	TypeList    FieldType = "list"
)

// FieldDef defines constraints for a single config field.
type FieldDef struct {
	Type     FieldType `yaml:"type"`
	Required bool      `yaml:"required"`
}

// Schema is a flat or nested map of field definitions.
type Schema map[string]FieldDef

// LoadSchema reads a YAML schema definition file and returns a Schema.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read file %q: %w", path, err)
	}
	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("schema: parse %q: %w", path, err)
	}
	return s, nil
}

// Validate checks that the provided config map conforms to the schema.
// It returns a list of human-readable violation messages.
func Validate(cfg map[string]any, s Schema) []string {
	var violations []string

	for key, def := range s {
		val, exists := cfg[key]
		if !exists {
			if def.Required {
				violations = append(violations, fmt.Sprintf("required key %q is missing", key))
			}
			continue
		}
		if def.Type != "" {
			if err := checkType(key, val, def.Type); err != nil {
				violations = append(violations, err.Error())
			}
		}
	}
	return violations
}

func checkType(key string, val any, expected FieldType) error {
	var ok bool
	switch expected {
	case TypeString:
		_, ok = val.(string)
	case TypeInt:
		switch val.(type) {
		case int, int64, float64:
			ok = true
		}
	case TypeBool:
		_, ok = val.(bool)
	case TypeFloat:
		switch val.(type) {
		case float32, float64:
			ok = true
		}
	case TypeMap:
		_, ok = val.(map[string]any)
	case TypeList:
		_, ok = val.([]any)
	default:
		return nil
	}
	if !ok {
		return fmt.Errorf("key %q: expected type %q, got %T", key, expected, val)
	}
	return nil
}
