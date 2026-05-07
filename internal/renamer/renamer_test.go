package renamer_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/renamer"
)

func baseMap() map[string]any {
	return map[string]any{
		"host": "localhost",
		"port": 5432,
		"database": map[string]any{
			"name": "mydb",
			"pool": map[string]any{
				"size": 10,
			},
		},
	}
}

func TestRename_TopLevelKey(t *testing.T) {
	result, err := renamer.Rename(baseMap(), []renamer.Rule{
		{OldPath: "host", NewPath: "hostname"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["host"]; ok {
		t.Error("old key 'host' should be removed")
	}
	if result["hostname"] != "localhost" {
		t.Errorf("expected hostname=localhost, got %v", result["hostname"])
	}
}

func TestRename_NestedKey(t *testing.T) {
	result, err := renamer.Rename(baseMap(), []renamer.Rule{
		{OldPath: "database.name", NewPath: "database.db_name"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, _ := result["database"].(map[string]any)
	if _, ok := db["name"]; ok {
		t.Error("old key 'database.name' should be removed")
	}
	if db["db_name"] != "mydb" {
		t.Errorf("expected db_name=mydb, got %v", db["db_name"])
	}
}

func TestRename_MoveKeyAcrossLevel(t *testing.T) {
	result, err := renamer.Rename(baseMap(), []renamer.Rule{
		{OldPath: "database.pool.size", NewPath: "pool_size"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["pool_size"] != 10 {
		t.Errorf("expected pool_size=10, got %v", result["pool_size"])
	}
}

func TestRename_MissingKey_ReturnsError(t *testing.T) {
	_, err := renamer.Rename(baseMap(), []renamer.Rule{
		{OldPath: "nonexistent", NewPath: "other"},
	})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	input := baseMap()
	_, err := renamer.Rename(input, []renamer.Rule{
		{OldPath: "host", NewPath: "hostname"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := input["host"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestRename_MultipleRules(t *testing.T) {
	result, err := renamer.Rename(baseMap(), []renamer.Rule{
		{OldPath: "host", NewPath: "hostname"},
		{OldPath: "port", NewPath: "db_port"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["hostname"] != "localhost" {
		t.Errorf("expected hostname=localhost, got %v", result["hostname"])
	}
	if result["db_port"] != 5432 {
		t.Errorf("expected db_port=5432, got %v", result["db_port"])
	}
}
