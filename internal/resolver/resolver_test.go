package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/confmerge/internal/resolver"
)

func TestResolveFiles_SingleFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.yaml")
	os.WriteFile(f, []byte("key: val\n"), 0644)

	layers, err := resolver.ResolveFiles([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(layers) != 1 {
		t.Fatalf("expected 1 layer, got %d", len(layers))
	}
	if layers[0].Path != f {
		t.Errorf("expected path %q, got %q", f, layers[0].Path)
	}
}

func TestResolveFiles_MultipleFiles_OrderPreserved(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a.yaml")
	b := filepath.Join(tmp, "b.toml")
	os.WriteFile(a, []byte("k: 1\n"), 0644)
	os.WriteFile(b, []byte("k = 2\n"), 0644)

	layers, err := resolver.ResolveFiles([]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(layers) != 2 {
		t.Fatalf("expected 2 layers, got %d", len(layers))
	}
	if layers[0].Priority >= layers[1].Priority {
		t.Errorf("expected layers[0].Priority < layers[1].Priority")
	}
}

func TestResolveFiles_Directory(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "00-base.yaml"), []byte("k: base\n"), 0644)
	os.WriteFile(filepath.Join(tmp, "01-override.yaml"), []byte("k: override\n"), 0644)
	os.WriteFile(filepath.Join(tmp, "ignored.txt"), []byte("nope"), 0644)

	layers, err := resolver.ResolveFiles([]string{tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(layers) != 2 {
		t.Fatalf("expected 2 layers, got %d", len(layers))
	}
}

func TestResolveFiles_UnsupportedExtension(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.ini")
	os.WriteFile(f, []byte("[section]\n"), 0644)

	_, err := resolver.ResolveFiles([]string{f})
	if err == nil {
		t.Fatal("expected error for unsupported extension, got nil")
	}
}

func TestResolveFiles_NotFound(t *testing.T) {
	_, err := resolver.ResolveFiles([]string{"/nonexistent/path/config.yaml"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
