package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Config holds the parsed CLI flags for a single confmerge invocation.
type Config struct {
	Inputs       []string
	Output       string
	Format       string
	Diff         bool
	DryRun       bool
	SortKeys     bool
	SortDesc     bool
	EnvSubst     bool
	Flatten      bool
	SchemaFile   string
	PatchFile    string
	Profile      bool
	CacheDir     string
	Watch        bool
}

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("confmerge", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		output     = fs.String("output", "", "write merged result to `file` (default: stdout)")
		format     = fs.String("format", "yaml", "output format: yaml|toml|json")
		diff       = fs.Bool("diff", false, "print a diff instead of the merged result")
		dryRun     = fs.Bool("dry-run", false, "validate and diff without writing output")
		sortKeys   = fs.Bool("sort-keys", false, "sort map keys in the output")
		sortDesc   = fs.Bool("sort-desc", false, "sort keys in descending order (requires --sort-keys)")
		envSubst   = fs.Bool("env-subst", false, "substitute ${VAR} placeholders from environment")
		flatten    = fs.Bool("flatten", false, "flatten nested maps to dot-separated keys")
		schema     = fs.String("schema", "", "path to JSON-schema YAML file for validation")
		patch      = fs.String("patch", "", "path to patch file (RFC-6902-style operations)")
		profile    = fs.Bool("profile", false, "print pipeline stage timings to stderr")
		cacheDir   = fs.String("cache-dir", "", "directory for caching parsed inputs")
		watch      = fs.Bool("watch", false, "re-run pipeline when input files change")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	inputs := fs.Args()
	if len(inputs) == 0 {
		return nil, errors.New("at least one input file or directory is required")
	}

	if err := validateOutputFormat(*format); err != nil {
		return nil, err
	}

	return &Config{
		Inputs:     inputs,
		Output:     *output,
		Format:     *format,
		Diff:       *diff,
		DryRun:     *dryRun,
		SortKeys:   *sortKeys,
		SortDesc:   *sortDesc,
		EnvSubst:   *envSubst,
		Flatten:    *flatten,
		SchemaFile: *schema,
		PatchFile:  *patch,
		Profile:    *profile,
		CacheDir:   *cacheDir,
		Watch:      *watch,
	}, nil
}

func validateOutputFormat(f string) error {
	switch f {
	case "yaml", "toml", "json":
		return nil
	}
	return fmt.Errorf("unsupported output format %q: must be yaml, toml, or json", f)
}
