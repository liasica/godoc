// Copyright (C) godoc. 2026-present.
//
// Created at 2026-01-12, by liasica

package godoc

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/labstack/echo/v4"
)

type HtmlTemplate struct {
	Templates map[string]*template.Template
}

func (t *HtmlTemplate) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.Templates[name].ExecuteTemplate(w, name, data)
}

// LoadTemplates Load HTML templates from embedded file system
func LoadTemplates(tmpls embed.FS, templatesDir string) (ht *HtmlTemplate) {
	ht = &HtmlTemplate{Templates: make(map[string]*template.Template)}

	_ = fs.WalkDir(tmpls, templatesDir, func(path string, d fs.DirEntry, _ error) (err error) {
		if d.IsDir() {
			return
		}

		name := strings.Replace(path, templatesDir+"/", "", 1)
		pt := template.New(name)
		b, _ := tmpls.ReadFile(path)
		_, _ = pt.Parse(string(b))

		ht.Templates[name] = pt
		return
	})

	return
}
