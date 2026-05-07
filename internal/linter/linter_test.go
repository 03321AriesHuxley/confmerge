package linter

import (
	"strings"
	"testing"
)

func TestLint_NoIssues(t *testing.T) {
	data := map[string]any{
		"host": "localhost",
		"port": 8080,
		"db": map[string]any{
			"name": "mydb",
		},
	}
	issues := Lint(data, DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestLint_DuplicateCaseInsensitiveKey(t *testing.T) {
	data := map[string]any{
		"Host": "localhost",
		"host": "remotehost",
	}
	issues := Lint(data, DefaultOptions())
	if len(issues) == 0 {
		t.Fatal("expected duplicate-key issue, got none")
	}
	if issues[0].Level != "error" {
		t.Errorf("expected level 'error', got %q", issues[0].Level)
	}
}

func TestLint_StringTooLong(t *testing.T) {
	opts := Options{MaxDepth: 10, MaxStringBytes: 10}
	data := map[string]any{
		"key": strings.Repeat("x", 20),
	}
	issues := Lint(data, opts)
	if len(issues) == 0 {
		t.Fatal("expected string-length issue, got none")
	}
	if issues[0].Level != "warn" {
		t.Errorf("expected level 'warn', got %q", issues[0].Level)
	}
}

func TestLint_ExceedsMaxDepth(t *testing.T) {
	opts := Options{MaxDepth: 2, MaxStringBytes: 0}
	data := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	issues := Lint(data, opts)
	if len(issues) == 0 {
		t.Fatal("expected depth issue, got none")
	}
	if !strings.Contains(issues[0].Message, "nesting depth") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_IssueStringFormat(t *testing.T) {
	issue := Issue{Level: "warn", Path: "db.host", Message: "too long"}
	s := issue.String()
	if !strings.Contains(s, "[WARN]") || !strings.Contains(s, "db.host") {
		t.Errorf("unexpected issue string: %s", s)
	}
}

func TestLint_NestedDuplicateKey(t *testing.T) {
	data := map[string]any{
		"db": map[string]any{
			"Name": "prod",
			"name": "staging",
		},
	}
	issues := Lint(data, DefaultOptions())
	if len(issues) == 0 {
		t.Fatal("expected nested duplicate-key issue")
	}
	if !strings.HasPrefix(issues[0].Path, "db.") {
		t.Errorf("expected path under 'db', got %q", issues[0].Path)
	}
}
