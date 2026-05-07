package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Config holds all parsed CLI flags.
type Config struct {
	Inputs       []string
	Output       string
	Format       string
	Diff         bool
	DryRun       bool
	FlattenKeys  bool
	EnvSubst     bool
	SchemaFile   string
	IncludeKeys  []string
	ExcludeKeys  []string
}

type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("confmerge", flag.ContinueOnError)

	var (
		output     = fs.String("output", "", "output file path (default: stdout)")
		format     = fs.String("format", "yaml", "output format: yaml|toml|json")
		diff       = fs.Bool("diff", false, "print a diff between first and merged result")
		dryRun     = fs.Bool("dry-run", false, "validate and print without writing output")
		flatten    = fs.Bool("flatten", false, "flatten nested keys using dot notation")
		envSubst   = fs.Bool("env-subst", false, "substitute ${ENV_VAR} references")
		schemaFile = fs.String("schema", "", "path to JSON schema file for validation")
		include    multiFlag
		exclude    multiFlag
	)

	fs.Var(&include, "include", "dot-path key to include (repeatable)")
	fs.Var(&exclude, "exclude", "dot-path key to exclude (repeatable)")

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
		Inputs:      inputs,
		Output:      *output,
		Format:      *format,
		Diff:        *diff,
		DryRun:      *dryRun,
		FlattenKeys: *flatten,
		EnvSubst:    *envSubst,
		SchemaFile:  *schemaFile,
		IncludeKeys: []string(include),
		ExcludeKeys: []string(exclude),
	}, nil
}

func validateOutputFormat(f string) error {
	switch f {
	case "yaml", "toml", "json":
		return nil
	}
	return fmt.Errorf("unsupported output format %q: must be yaml, toml, or json", f)
}
