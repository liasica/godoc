// Copyright (C)  2026-present.
//
// Created at 2026-01-12, by liasica

package godoc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/swaggo/swag/v2/format"
	"github.com/swaggo/swag/v2/gen"
)

// Generate generates documentation based on the provided configuration file path.
func Generate(cfgPath string) error {
	basePath := filepath.Dir(cfgPath)

	// parse config
	cfg, err := LoadConfig(ResolveConfigPath(cfgPath))
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v\n", err)
	}

	var gim *GoMod
	gim, err = NewGoMod(basePath)
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %v\n", err)
	}

	externalPaths := cfg.ExternalPath
	paths := make([]string, len(cfg.Path))
	for i, cp := range cfg.Path {
		paths[i] = filepath.Join(basePath, cp)
	}

	for ep, sub := range externalPaths {
		var p string
		p, err = gim.GetPath(basePath, ep, sub)
		if err != nil {
			return fmt.Errorf("dependency resolution failed: %v\n", err)
		}
		paths = append(paths, p)
	}

	mainFile := cfg.MainFile
	output := filepath.Join(basePath, cfg.Output)
	searchDir := strings.Join(paths, ",")

	fmt.Printf("starting documentation generation: main=%s, deps=%s, output=%s\n", mainFile, searchDir, output)

	fmt.Println("formatting...")

	err = format.New().Build(&format.Config{
		SearchDir: searchDir,
		MainFile:  mainFile,
	})

	if err != nil {
		return fmt.Errorf("formatting failed: %v\n", err)
	}

	fmt.Println("generating documentation...")
	gc := &gen.Config{
		SearchDir:   searchDir,
		MainAPIFile: mainFile,
		// ParseDependency: 1,
		OutputDir:           output,
		OutputTypes:         cfg.OutputTypes,
		GenerateOpenAPI3Doc: true,
	}

	// resolve Markdown files directory
	if cfg.MarkdownFilesDir != "" {
		gc.MarkdownFilesDir, err = cfg.ResolveAbsPath(cfg.MarkdownFilesDir)
		if err != nil {
			return fmt.Errorf("failed to resolve markdown files directory: %v\n", err)
		}
	}

	err = gen.New().Build(gc)
	if err != nil {
		return fmt.Errorf("documentation generation failed: %v\n", err)
	}

	err = ConvertEnum2OneOf(filepath.Join(output, "swagger.yaml"))
	if err != nil {
		return fmt.Errorf("enum conversion failed: %v\n", err)
	}

	return nil
}
