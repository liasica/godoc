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

func Preview() (*cobra.Group, *cobra.Command) {
	g := &cobra.Group{
		ID:    "preview",
		Title: "Documentation Preview Commands",
	}

	var (
		address  string
		generate bool
	)

	cmd := &cobra.Command{
		Use:               "preview",
		Short:             "Preview documentation for Go projects",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		GroupID:           g.ID,
		Run: func(_ *cobra.Command, _ []string) {
			var err error
			if generate {
				err = godoc.Generate(cfgPath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			// generate a random available port if address is empty
			if address == "" {
				address, err = godoc.GetAvailableAddress()
				if err != nil {
					fmt.Printf("failed to get available address: %v\n", err)
					os.Exit(1)
				}
			}

			err = godoc.Preview(cfgPath, address)
			if err != nil {
				fmt.Printf("failed to start preview server: %v\n", err)
				os.Exit(1)
			}

		},
	}

	cmd.Flags().StringVarP(&address, "address", "a", "", "address to run the preview server on, if empty, a random available port will be used")
	cmd.Flags().BoolVarP(&generate, "generate", "g", false, "generate documentation before starting the preview server")

	return g, cmd
}
