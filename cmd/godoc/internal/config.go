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

func Config() (*cobra.Group, *cobra.Command) {
	g := &cobra.Group{
		ID:    "config",
		Title: "Configuration Commands",
	}

	cmd := &cobra.Command{
		Use:               "config",
		Short:             "Manage configuration",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		GroupID:           g.ID,
	}

	cmd.AddCommand(configInit())

	return g, cmd
}

func configInit() *cobra.Command {
	return &cobra.Command{
		Use:               "init",
		Short:             "Initialize configuration",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			// detect if config file already exists
			if _, err := os.Stat(".godoc.yaml"); err == nil {
				fmt.Println("configuration file .godoc.yaml already exists")
				os.Exit(1)
			}

			// write default config to .godoc.yaml
			defaultCfg := godoc.DefaultConfig()
			fmt.Printf("default configuration:\n%s\n", defaultCfg)
			err := os.WriteFile(".godoc.yaml", []byte(defaultCfg), os.ModePerm)
			if err != nil {
				fmt.Printf("failed to write configuration file: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("configuration file .godoc.yaml created successfully")
		},
	}
}
