package main

import (
	"errors"
	"flag"
	"fmt"
)

// config holds the parsed CLI configuration.
type config struct {
	inputs       []string
	output       string
	outputFormat string
}

// parseFlags parses the command-line arguments and returns a config.
func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("confmerge", flag.ContinueOnError)

	var output string
	var outputFormat string

	fs.StringVar(&output, "o", "", "Output file path (default: stdout)")
	fs.StringVar(&output, "output", "", "Output file path (default: stdout)")
	fs.StringVar(&outputFormat, "f", "", "Output format: yaml, toml, json (auto-detected from -o if omitted)")
	fs.StringVar(&outputFormat, "format", "", "Output format: yaml, toml, json (auto-detected from -o if omitted)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: confmerge [options] <file-or-dir>...\n\nOptions:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	inputs := fs.Args()
	if len(inputs) == 0 {
		return nil, errors.New("at least one input file or directory is required")
	}

	if err := validateOutputFormat(outputFormat); err != nil {
		return nil, err
	}

	return &config{
		inputs:       inputs,
		output:       output,
		outputFormat: outputFormat,
	}, nil
}

// validateOutputFormat checks that the given format is one of the supported
// output formats. An empty string is accepted and means auto-detection.
func validateOutputFormat(format string) error {
	if format == "" {
		return nil
	}
	switch format {
	case "yaml", "toml", "json":
		return nil
	default:
		return fmt.Errorf("unsupported output format %q: must be yaml, toml, or json", format)
	}
}
