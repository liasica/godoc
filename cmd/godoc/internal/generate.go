// Copyright (C) godoc. 2026-present.
//
// Created at 2026-01-12, by liasica

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/swaggo/swag/v2/format"
	"github.com/swaggo/swag/v2/gen"

	"github.com/liasica/godoc"
)

func Generate() (*cobra.Group, *cobra.Command) {
	g := &cobra.Group{
		ID:    "generate",
		Title: "Documentation Generation Commands",
	}

	var cfgPath string

	cmd := &cobra.Command{
		Use:               "generate",
		Short:             "Generate documentation for Go projects",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		GroupID:           g.ID,
		Run: func(_ *cobra.Command, _ []string) {
			// parse config
			cfg, err := godoc.LoadConfig(godoc.ResolveConfigPath(cfgPath))
			if err != nil {
				fmt.Printf("failed to load configuration: %v\n", err)
				os.Exit(1)
			}

			var gim *godoc.GoMod
			gim, err = godoc.NewGoMod()
			if err != nil {
				fmt.Printf("dependency resolution failed: %v\n", err)
				os.Exit(1)
			}

			externalPaths := cfg.ExternalPath
			paths := cfg.Path

			for ep, sub := range externalPaths {
				var p string
				p, err = gim.GetPath(ep, sub)
				if err != nil {
					fmt.Printf("dependency resolution failed: %v\n", err)
					os.Exit(1)
				}
				paths = append(paths, p)
			}

			mainFile := cfg.MainFile
			output := cfg.Output
			searchDir := strings.Join(paths, ",")

			fmt.Printf("starting documentation generation: main=%s, deps=%s, output=%s", mainFile, searchDir, output)

			fmt.Println("formatting...")
			err = format.New().Build(&format.Config{
				SearchDir: searchDir,
				MainFile:  mainFile,
			})
			if err != nil {
				fmt.Printf("formatting failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("generating documentation...")
			err = gen.New().Build(&gen.Config{
				SearchDir:   searchDir,
				MainAPIFile: mainFile,
				// ParseDependency: 1,
				OutputDir:           output,
				OutputTypes:         cfg.OutputTypes,
				GenerateOpenAPI3Doc: true,
			})
			if err != nil {
				fmt.Printf("documentation generation failed: %v\n", err)
				os.Exit(1)
			}

			err = godoc.ConvertEnum2OneOf(filepath.Join(output, "swagger.yaml"))
			if err != nil {
				fmt.Printf("enum conversion failed: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&cfgPath, "config", "c", ".godoc.yaml", "path to config YAML file")
	return g, cmd
}
