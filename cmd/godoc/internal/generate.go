// Copyright (C) godoc. 2026-present.
//
// Created at 2026-01-12, by liasica

package internal

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

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
			err := godoc.Generate(cfgPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&cfgPath, "config", "c", ".godoc.yaml", "path to config YAML file")
	return g, cmd
}
