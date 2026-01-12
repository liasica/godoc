// Copyright (C) moneta. 2025-present.
//
// Created at 2025-08-04, by liasica

package assets

import (
	"embed"
	_ "embed"
)

var (
	//go:embed templates/*
	TemplateFS embed.FS
)
