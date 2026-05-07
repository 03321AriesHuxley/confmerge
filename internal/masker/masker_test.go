package masker_test

import (
	"strings"
	"testing"

	"github.com/yourorg/confmerge/internal/masker"
)

func TestMask_ExactValue(t *testing.T) {
	input := map[string]any{"token": "abc123", "name": "alice"}
	out := masker.Mask(input, masker.Options{Values: []string{"abc123"}})
	if out["token"] != masker.DefaultPlaceholder {
		t.Errorf("expected token to be masked, got %v", out["token"])
	}
	if out["name"] != "alice" {
		t.Errorf("expected name to be unchanged, got %v", out["name"])
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	input := map[string]any{"secret": "hunter2"}
	out := masker.Mask(input, masker.Options{Values: []string{"hunter2"}, Placeholder: "<hidden>"})
	if out["secret"] != "<hidden>" {
		t.Errorf("expected <hidden>, got %v", out["secret"])
	}
}

func TestMask_MatchFn(t *testing.T) {
	input := map[string]any{"api_key": "xyz", "host": "localhost"}
	opts := masker.Options{
		MatchFn: func(key, _ string) bool {
			return strings.HasSuffix(key, "_key")
		},
	}
	out := masker.Mask(input, opts)
	if out["api_key"] != masker.DefaultPlaceholder {
		t.Errorf("expected api_key masked")
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host unchanged")
	}
}

func TestMask_NestedMap(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{"password": "s3cr3t", "user": "admin"},
	}
	out := masker.Mask(input, masker.Options{Values: []string{"s3cr3t"}})
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	if db["password"] != masker.DefaultPlaceholder {
		t.Errorf("expected nested password masked")
	}
	if db["user"] != "admin" {
		t.Errorf("expected user unchanged")
	}
}

func TestMask_SliceValues(t *testing.T) {
	input := map[string]any{"tokens": []any{"tok1", "tok2", "safe"}}
	out := masker.Mask(input, masker.Options{Values: []string{"tok1", "tok2"}})
	slice, ok := out["tokens"].([]any)
	if !ok || len(slice) != 3 {
		t.Fatal("expected slice of 3")
	}
	if slice[0] != masker.DefaultPlaceholder || slice[1] != masker.DefaultPlaceholder {
		t.Errorf("expected first two elements masked")
	}
	if slice[2] != "safe" {
		t.Errorf("expected third element unchanged")
	}
}

func TestMask_DoesNotMutateInput(t *testing.T) {
	input := map[string]any{"pw": "original"}
	masker.Mask(input, masker.Options{Values: []string{"original"}})
	if input["pw"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestMask_NonStringValuesUnchanged(t *testing.T) {
	input := map[string]any{"port": 8080, "enabled": true}
	out := masker.Mask(input, masker.Options{Values: []string{"8080"}})
	if out["port"] != 8080 {
		t.Errorf("expected int unchanged, got %v", out["port"])
	}
}
