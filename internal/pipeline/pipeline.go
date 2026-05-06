// Package pipeline orchestrates the full confmerge processing pipeline,
// wiring together the resolver, loader, merger, transformer, validator,
// differ, and writer stages into a single cohesive execution flow.
package pipeline

import (
	"fmt"
	"io"

	"github.com/yourorg/confmerge/internal/differ"
	"github.com/yourorg/confmerge/internal/loader"
	"github.com/yourorg/confmerge/internal/merger"
	"github.com/yourorg/confmerge/internal/profiler"
	"github.com/yourorg/confmerge/internal/resolver"
	"github.com/yourorg/confmerge/internal/schema"
	"github.com/yourorg/confmerge/internal/transformer"
	"github.com/yourorg/confmerge/internal/validator"
	"github.com/yourorg/confmerge/internal/writer"
)

// Options holds all configuration for a pipeline run.
type Options struct {
	// Inputs is a list of file paths or directories to merge, in order.
	Inputs []string

	// OutputPath is the destination file. If empty, output goes to Stdout.
	OutputPath string

	// OutputFormat overrides the output format (yaml, toml, json).
	// When empty it is inferred from OutputPath or defaults to "yaml".
	OutputFormat string

	// SchemaPath is an optional JSON-Schema file used to validate the result.
	SchemaPath string

	// TransformOpts are forwarded to the transformer stage.
	TransformOpts transformer.Options

	// Diff, when true, prints a human-readable diff to DiffWriter instead of
	// writing the merged output.
	Diff bool

	// DiffWriter is where diff output is written. Defaults to os.Stdout.
	DiffWriter io.Writer

	// Profile, when true, prints timing information after the run.
	Profile bool

	// ProfileWriter is where profiling output is written. Defaults to os.Stderr.
	ProfileWriter io.Writer
}

// Run executes the full merge pipeline according to opts and writes the result
// to opts.OutputPath (or opts.DiffWriter when opts.Diff is true).
// It returns an error if any stage fails.
func Run(opts Options) error {
	p := profiler.New()

	// ── 1. Resolve inputs ────────────────────────────────────────────────────
	p.Track("resolve", func() {})
	files, err := resolver.ResolveFiles(opts.Inputs)
	if err != nil {
		return fmt.Errorf("resolve: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("resolve: no supported config files found in provided inputs")
	}

	// ── 2. Load & merge layers ────────────────────────────────────────────────
	var merged map[string]any
	for _, f := range files {
		var layerData map[string]any
		p.Track("load:"+f, func() {})
		layerData, err = loader.LoadFile(f)
		if err != nil {
			return fmt.Errorf("load %q: %w", f, err)
		}
		p.Track("merge:"+f, func() {})
		merged = merger.Merge(merged, layerData)
	}

	// ── 3. Transform ─────────────────────────────────────────────────────────
	p.Track("transform", func() {})
	merged, err = transformer.Apply(merged, opts.TransformOpts)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	// ── 4. Validate (structural) ──────────────────────────────────────────────
	p.Track("validate", func() {})
	if errs := validator.Validate(merged); len(errs) > 0 {
		return fmt.Errorf("validate: %w", errs[0])
	}

	// ── 5. Schema validation (optional) ──────────────────────────────────────
	if opts.SchemaPath != "" {
		p.Track("schema", func() {})
		sc, loadErr := schema.LoadSchema(opts.SchemaPath)
		if loadErr != nil {
			return fmt.Errorf("schema load: %w", loadErr)
		}
		if schErr := sc.Validate(merged); schErr != nil {
			return fmt.Errorf("schema validate: %w", schErr)
		}
	}

	// ── 6. Diff mode ─────────────────────────────────────────────────────────
	if opts.Diff {
		p.Track("diff", func() {})
		// Load the first file as the "base" and diff against the fully merged result.
		base, loadErr := loader.LoadFile(files[0])
		if loadErr != nil {
			return fmt.Errorf("diff base load: %w", loadErr)
		}
		changes := differ.Diff(base, merged)
		fmt.Fprint(opts.DiffWriter, differ.FormatText(changes))
		fmt.Fprintf(opts.DiffWriter, "\n%s\n", differ.Summary(changes))
		if opts.Profile {
			p.Print(opts.ProfileWriter)
		}
		return nil
	}

	// ── 7. Write output ───────────────────────────────────────────────────────
	p.Track("write", func() {})
	fmt := opts.OutputFormat
	if fmt == "" {
		fmt = "yaml"
	}
	if err = writer.Write(merged, opts.OutputPath, fmt); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	if opts.Profile {
		p.Print(opts.ProfileWriter)
	}
	return nil
}
