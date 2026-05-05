package main

import (
	"testing"
)

func TestParseFlags_MinimalArgs(t *testing.T) {
	cfg, err := parseFlags([]string{"base.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.inputs) != 1 || cfg.inputs[0] != "base.yaml" {
		t.Errorf("expected inputs [base.yaml], got %v", cfg.inputs)
	}
	if cfg.output != "" {
		t.Errorf("expected empty output, got %q", cfg.output)
	}
}

func TestParseFlags_WithOutputAndFormat(t *testing.T) {
	cfg, err := parseFlags([]string{"-o", "out.json", "-f", "json", "a.yaml", "b.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.output != "out.json" {
		t.Errorf("expected output out.json, got %q", cfg.output)
	}
	if cfg.outputFormat != "json" {
		t.Errorf("expected format json, got %q", cfg.outputFormat)
	}
	if len(cfg.inputs) != 2 {
		t.Errorf("expected 2 inputs, got %d", len(cfg.inputs))
	}
}

func TestParseFlags_NoInputs(t *testing.T) {
	_, err := parseFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing inputs, got nil")
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := parseFlags([]string{"-f", "xml", "base.yaml"})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestParseFlags_MultipleInputs(t *testing.T) {
	cfg, err := parseFlags([]string{"base.yaml", "override/", "local.toml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.inputs) != 3 {
		t.Errorf("expected 3 inputs, got %d", len(cfg.inputs))
	}
}
