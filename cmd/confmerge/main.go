// Package main is the entry point for the confmerge CLI tool.
package main

import (
	"fmt"
	"os"

	"github.com/yourorg/confmerge/internal/loader"
	"github.com/yourorg/confmerge/internal/merger"
	"github.com/yourorg/confmerge/internal/resolver"
	"github.com/yourorg/confmerge/internal/writer"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	files, err := resolver.ResolveFiles(cfg.inputs)
	if err != nil {
		return fmt.Errorf("resolving inputs: %w", err)
	}

	base := map[string]interface{}{}
	for _, f := range files {
		data, err := loader.LoadFile(f)
		if err != nil {
			return fmt.Errorf("loading %s: %w", f, err)
		}
		base = merger.Merge(base, data)
	}

	format := cfg.outputFormat
	if format == "" {
		format = loader.DetectFormat(cfg.output)
		if format == "" {
			format = "yaml"
		}
	}

	var out *os.File
	if cfg.output == "" || cfg.output == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(cfg.output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer out.Close()
	}

	return writer.Write(out, base, format)
}
