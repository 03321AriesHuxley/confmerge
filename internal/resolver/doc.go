// Package resolver handles discovery and ordering of configuration layer files
// for confmerge.
//
// Given a list of explicit file paths and/or directory paths, the resolver
// determines which files are valid configuration sources (YAML, TOML, JSON)
// and returns them sorted by their merge precedence (lowest priority first,
// highest priority last).
//
// When a directory is provided, all supported files within it are included and
// ordered alphabetically, making it easy to use numeric prefixes (e.g.
// 00-base.yaml, 10-env.yaml) to control layer ordering.
//
// Explicit file paths maintain the order they were supplied on the command line,
// with each path receiving a higher base priority than the previous one.
package resolver
