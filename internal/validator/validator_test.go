package validator

import (
	"strings"
	"testing"
)

func TestValidate_ValidMap(t *testing.T) {
	data := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
		"app": map[string]interface{}{
			"name": "myapp",
		},
	}
	if err := Validate(data); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidate_NullValue(t *testing.T) {
	data := map[string]interface{}{
		"key": nil,
	}
	err := Validate(data)
	if err == nil {
		t.Fatal("expected error for null value, got nil")
	}
	if !strings.Contains(err.Error(), "null value") {
		t.Errorf("expected 'null value' in error, got: %v", err)
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	data := map[string]interface{}{
		"": "oops",
	}
	err := Validate(data)
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
	if !strings.Contains(err.Error(), "empty key") {
		t.Errorf("expected 'empty key' in error, got: %v", err)
	}
}

func TestValidate_NestedNullValue(t *testing.T) {
	data := map[string]interface{}{
		"server": map[string]interface{}{
			"host": nil,
		},
	}
	err := Validate(data)
	if err == nil {
		t.Fatal("expected error for nested null value, got nil")
	}
	if !strings.Contains(err.Error(), "server.host") {
		t.Errorf("expected path 'server.host' in error, got: %v", err)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	data := map[string]interface{}{
		"a": nil,
		"b": nil,
	}
	err := Validate(data)
	if err == nil {
		t.Fatal("expected errors, got nil")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(ve.Errors))
	}
}

func TestValidate_EmptyMap(t *testing.T) {
	data := map[string]interface{}{}
	if err := Validate(data); err != nil {
		t.Errorf("expected no error for empty map, got: %v", err)
	}
}
