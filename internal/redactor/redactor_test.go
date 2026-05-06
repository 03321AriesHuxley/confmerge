package redactor_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/redactor"
)

func TestRedact_DefaultPatterns(t *testing.T) {
	input := map[string]any{
		"host":     "localhost",
		"password": "s3cr3t",
		"port":     5432,
	}
	out := redactor.Redact(input, redactor.Options{})
	if out["password"] != "***" {
		t.Errorf("expected password to be masked, got %v", out["password"])
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host to be unchanged, got %v", out["host"])
	}
	if out["port"] != 5432 {
		t.Errorf("expected port to be unchanged, got %v", out["port"])
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	input := map[string]any{
		"db_pass": "hunter2",
		"my_pin":  "1234",
		"name":    "alice",
	}
	opts := redactor.Options{Patterns: []string{"pin"}, Mask: "[HIDDEN]"}
	out := redactor.Redact(input, opts)
	if out["my_pin"] != "[HIDDEN]" {
		t.Errorf("expected my_pin masked, got %v", out["my_pin"])
	}
	if out["db_pass"] == "[HIDDEN]" {
		t.Error("db_pass should not be masked with custom patterns")
	}
	if out["name"] != "alice" {
		t.Errorf("expected name unchanged, got %v", out["name"])
	}
}

func TestRedact_NestedMap(t *testing.T) {
	input := map[string]any{
		"database": map[string]any{
			"host":   "db.example.com",
			"secret": "topsecret",
		},
		"app": "myapp",
	}
	out := redactor.Redact(input, redactor.Options{})
	db, ok := out["database"].(map[string]any)
	if !ok {
		t.Fatal("expected database to be a map")
	}
	if db["secret"] != "***" {
		t.Errorf("expected nested secret masked, got %v", db["secret"])
	}
	if db["host"] != "db.example.com" {
		t.Errorf("expected nested host unchanged, got %v", db["host"])
	}
}

func TestRedact_CaseInsensitiveKey(t *testing.T) {
	input := map[string]any{
		"API_KEY": "abc123",
		"value":   42,
	}
	out := redactor.Redact(input, redactor.Options{})
	if out["API_KEY"] != "***" {
		t.Errorf("expected API_KEY masked, got %v", out["API_KEY"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	input := map[string]any{
		"token": "original",
	}
	redactor.Redact(input, redactor.Options{})
	if input["token"] != "original" {
		t.Error("Redact must not mutate the input map")
	}
}
