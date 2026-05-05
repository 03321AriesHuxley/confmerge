package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LayerFile represents a single config layer with its precedence order.
type LayerFile struct {
	Path     string
	Priority int
}

// ResolveFiles takes a list of file paths or a directory and returns an ordered
// slice of LayerFile values sorted by ascending priority (lowest = base layer).
func ResolveFiles(inputs []string) ([]LayerFile, error) {
	var layers []LayerFile

	for i, input := range inputs {
		info, err := os.Stat(input)
		if err != nil {
			return nil, fmt.Errorf("resolver: cannot access %q: %w", input, err)
		}

		if info.IsDir() {
			dirLayers, err := resolveDir(input, i*1000)
			if err != nil {
				return nil, err
			}
			layers = append(layers, dirLayers...)
		} else {
			if !isSupportedExt(input) {
				return nil, fmt.Errorf("resolver: unsupported file extension for %q", input)
			}
			layers = append(layers, LayerFile{Path: input, Priority: i})
		}
	}

	sort.Slice(layers, func(i, j int) bool {
		return layers[i].Priority < layers[j].Priority
	})

	return layers, nil
}

// resolveDir scans a directory for supported config files and assigns priorities
// based on alphabetical order within the directory, offset by baseOffset.
func resolveDir(dir string, baseOffset int) ([]LayerFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("resolver: cannot read directory %q: %w", dir, err)
	}

	var layers []LayerFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !isSupportedExt(name) {
			continue
		}
		layers = append(layers, LayerFile{
			Path:     filepath.Join(dir, name),
			Priority: baseOffset + len(layers),
		})
	}
	return layers, nil
}

func isSupportedExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml" || ext == ".toml" || ext == ".json"
}
