package transformer

import (
	"os"
	"testing"
)

func TestApply_NoOpts(t *testing.T) {
	input := map[string]any{"key": "value", "num": 42}
	out, err := Apply(input, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "value" || out["num"] != 42 {
		t.Errorf("expected unchanged map, got %v", out)
	}
}

func TestApply_EnvSubst(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	input := map[string]any{"env": "${APP_ENV}", "other": 1}
	out, err := Apply(input, Options{EnvSubst: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["env"] != "production" {
		t.Errorf("expected 'production', got %v", out["env"])
	}
}

func TestApply_EnvSubst_Nested(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	t.Cleanup(func() { os.Unsetenv("DB_HOST") })
	input := map[string]any{
		"database": map[string]any{
			"host": "$DB_HOST",
			"port": 5432,
		},
	}
	out, err := Apply(input, Options{EnvSubst: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["database"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected 'localhost', got %v", db["host"])
	}
}

func TestApply_Flatten(t *testing.T) {
	input := map[string]any{
		"server": map[string]any{
			"host": "0.0.0.0",
			"port": 8080,
		},
		"debug": true,
	}
	out, err := Apply(input, Options{FlattenSep: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["server.host"] != "0.0.0.0" {
		t.Errorf("expected 'server.host' key, got %v", out)
	}
	if out["server.port"] != 8080 {
		t.Errorf("expected 'server.port' key, got %v", out)
	}
	if out["debug"] != true {
		t.Errorf("expected 'debug' key, got %v", out)
	}
}

func TestApply_FlattenAndEnvSubst(t *testing.T) {
	t.Setenv("LOG_LEVEL", "info")
	input := map[string]any{
		"logging": map[string]any{
			"level": "${LOG_LEVEL}",
		},
	}
	out, err := Apply(input, Options{EnvSubst: true, FlattenSep: "_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["logging_level"] != "info" {
		t.Errorf("expected 'info', got %v", out["logging_level"])
	}
}
