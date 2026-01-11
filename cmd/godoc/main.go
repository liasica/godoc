// Copyright (C) moneta. 2025-present.
//
// Created at 2025-08-01, by liasica

package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/swaggo/swag/v2/format"
	"github.com/swaggo/swag/v2/gen"

	"github.com/liasica/godoc"
)

func main() {
	// --version
	versionFlag := flag.Bool("version", false, "print version information")
	// shorthand -v for version
	shortVersion := flag.Bool("v", false, "shorthand for --version")

	// --config
	cfgPath := flag.String("config", ".godoc.yaml", "path to config YAML file")
	flag.Parse()

	if *versionFlag || *shortVersion {
		fmt.Println(godoc.FullVersion())
		return
	}

	cfg, err := godoc.LoadConfig(godoc.ResolveConfigPath(*cfgPath))
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	var gim *godoc.GoMod
	gim, err = godoc.NewGoMod()
	if err != nil {
		log.Fatalf("dependency resolution failed: %v", err)
	}

	externalPaths := cfg.ExternalPath
	paths := cfg.Path

	for ep, sub := range externalPaths {
		var p string
		p, err = gim.GetPath(ep, sub)
		if err != nil {
			log.Fatalf("dependency resolution failed: %v", err)
		}
		paths = append(paths, p)
	}

	mainFile := cfg.MainFile
	output := cfg.Output
	searchDir := strings.Join(paths, ",")

	log.Printf("starting documentation generation: main=%s, deps=%s, output=%s", mainFile, searchDir, output)

	log.Println("formatting...")
	err = format.New().Build(&format.Config{
		SearchDir: searchDir,
		MainFile:  mainFile,
	})
	if err != nil {
		log.Fatalf("formatting failed: %v", err)
	}

	log.Println("generating documentation...")
	err = gen.New().Build(&gen.Config{
		SearchDir:   searchDir,
		MainAPIFile: mainFile,
		// ParseDependency: 1,
		OutputDir:           output,
		OutputTypes:         cfg.OutputTypes,
		GenerateOpenAPI3Doc: true,
	})
	if err != nil {
		log.Fatalf("documentation generation failed: %v", err)
	}

	err = godoc.ConvertEnum2OneOf(filepath.Join(output, "swagger.yaml"))
	if err != nil {
		log.Fatalf("enum conversion failed: %v", err)
	}
}
