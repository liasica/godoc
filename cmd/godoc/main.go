// Copyright (C) moneta. 2025-present.
//
// Created at 2025-08-01, by liasica

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/liasica/godoc"
	"github.com/liasica/godoc/cmd/godoc/internal"
)

func main() {
	cmd := cobra.Command{
		Use:               "godoc",
		Short:             "godoc is a documentation generator using swaggo/swag for Go projects",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Version:           godoc.GetVersion(),
	}

	genGroup, genCommand := internal.Generate()
	cfgGroup, cfgCommand := internal.Config()
	previewGroup, previewCommand := internal.Preview()

	cmd.AddGroup(
		genGroup,
		cfgGroup,
		previewGroup,
	)

	cmd.AddCommand(
		genCommand,
		cfgCommand,
		previewCommand,
	)

	err := cmd.Execute()
	if err != nil {
		fmt.Printf("command execution failed: %v\n", err)
		os.Exit(1)
	}
}
