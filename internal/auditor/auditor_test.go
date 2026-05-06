package auditor

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew_DefaultsToStderr(t *testing.T) {
	a := New(nil)
	if a == nil {
		t.Fatal("expected non-nil auditor")
	}
}

func TestLog_RecordsEntry(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf)
	a.Log("INFO", "merge", "started")

	if a.Count() != 1 {
		t.Fatalf("expected 1 entry, got %d", a.Count())
	}
	e := a.Entries()[0]
	if e.Level != "INFO" {
		t.Errorf("expected level INFO, got %s", e.Level)
	}
	if e.Stage != "merge" {
		t.Errorf("expected stage merge, got %s", e.Stage)
	}
	if e.Message != "started" {
		t.Errorf("expected message 'started', got %s", e.Message)
	}
}

func TestInfo_WritesToOutput(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf)
	a.Info("load", "files loaded successfully")

	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "load") {
		t.Errorf("expected stage 'load' in output, got: %s", out)
	}
}

func TestWarn_RecordsWarnLevel(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf)
	a.Warn("validate", "nullable field detected")

	if a.Count() != 1 {
		t.Fatalf("expected 1 entry, got %d", a.Count())
	}
	if a.Entries()[0].Level != "WARN" {
		t.Errorf("expected WARN level")
	}
}

func TestError_RecordsErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf)
	a.Error("schema", "validation failed")

	if a.Entries()[0].Level != "ERROR" {
		t.Errorf("expected ERROR level")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf)
	a.Info("stage1", "msg1")
	a.Info("stage2", "msg2")

	entries := a.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// Mutate copy — should not affect internal state
	entries[0].Message = "tampered"
	if a.Entries()[0].Message == "tampered" {
		t.Error("Entries should return a copy, not a reference")
	}
}

func TestMultipleLevels_OrderPreserved(t *testing.T) {
	var buf bytes.Buffer
	a := New(&buf)
	a.Info("a", "first")
	a.Warn("b", "second")
	a.Error("c", "third")

	entries := a.Entries()
	if entries[0].Stage != "a" || entries[1].Stage != "b" || entries[2].Stage != "c" {
		t.Error("entries not in insertion order")
	}
}
