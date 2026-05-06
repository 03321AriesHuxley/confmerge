package patcher_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/patcher"
)

func baseMap() map[string]interface{} {
	return map[string]interface{}{
		"app": map[string]interface{}{
			"name": "myapp",
			"port": 8080,
		},
		"debug": false,
	}
}

func TestApply_SetScalar(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpSet, Path: "debug", Value: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["debug"] != true {
		t.Errorf("expected debug=true, got %v", cfg["debug"])
	}
}

func TestApply_SetNestedPath(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpSet, Path: "app.port", Value: 9090},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := cfg["app"].(map[string]interface{})
	if app["port"] != 9090 {
		t.Errorf("expected port=9090, got %v", app["port"])
	}
}

func TestApply_SetCreatesIntermediateKeys(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpSet, Path: "db.host", Value: "localhost"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := cfg["db"].(map[string]interface{})
	if !ok {
		t.Fatal("expected db to be a map")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %v", db["host"])
	}
}

func TestApply_DeleteKey(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpDelete, Path: "app.port"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := cfg["app"].(map[string]interface{})
	if _, exists := app["port"]; exists {
		t.Error("expected port to be deleted")
	}
}

func TestApply_DeleteMissingKeyIsNoop(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpDelete, Path: "nonexistent.key"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_MergeAtRoot(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpMerge, Path: "", Value: map[string]interface{}{"env": "prod"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", cfg["env"])
	}
}

func TestApply_MergeNestedMap(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpMerge, Path: "app", Value: map[string]interface{}{"version": "1.2", "port": 443}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := cfg["app"].(map[string]interface{})
	if app["version"] != "1.2" {
		t.Errorf("expected version=1.2, got %v", app["version"])
	}
	if app["port"] != 443 {
		t.Errorf("expected port=443, got %v", app["port"])
	}
	if app["name"] != "myapp" {
		t.Errorf("expected name=myapp to be preserved, got %v", app["name"])
	}
}

func TestApply_UnknownOpReturnsError(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: "upsert", Path: "debug", Value: true},
	})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestApply_MergeNonMapValueReturnsError(t *testing.T) {
	cfg := baseMap()
	_, err := patcher.Apply(cfg, []patcher.Patch{
		{Op: patcher.OpMerge, Path: "debug", Value: "not-a-map"},
	})
	if err == nil {
		t.Error("expected error when merge value is not a map")
	}
}
