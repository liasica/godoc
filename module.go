// Copyright (C) moneta. 2025-present.
//
// Created at 2025-08-02, by liasica

package godoc

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/mod/modfile"
)

type GoMod struct {
	modfile      *modfile.File
	modcachePath string
}

func NewGoMod(basePath string) (gm *GoMod, err error) {
	p := filepath.Join(basePath, "go.mod")

	var b []byte
	b, err = os.ReadFile(p)
	if err != nil {
		return
	}

	var f *modfile.File
	f, err = modfile.Parse(p, b, nil)
	if err != nil {
		return
	}

	cache := GetGomodcache()
	if cache == "" {
		return nil, errors.New("failed to read GOMODCACHE environment variable")
	}

	return &GoMod{
		modfile:      f,
		modcachePath: cache,
	}, nil
}

func (gm *GoMod) GetPath(basePath string, module string, subpath string) (p string, err error) {
	var found *modfile.Require
	for _, require := range gm.modfile.Require {
		if require.Mod.Path == module {
			found = require
		}
	}

	if found == nil {
		return "", errors.New("module not found: " + module)
	}

	p = filepath.Join(gm.modcachePath, found.Mod.Path+"@"+found.Mod.Version)

	for _, replace := range gm.modfile.Replace {
		if replace.Old.Path == module {
			if isLocalPath(replace.New.Path) {
				// If it's a local path, return it directly
				p = filepath.Join(basePath, replace.New.Path)
			} else {
				// Otherwise return the replaced path inside the module cache
				p = filepath.Join(gm.modcachePath, replace.New.Path+"@"+replace.New.Version)
			}
		}
	}

	if subpath != "" {
		p = filepath.Join(p, subpath)
	}

	if _, err = os.Stat(p); err != nil {
		return "", errors.New("module path does not exist or is inaccessible: " + p)
	}

	return
}

// Determine whether this is a local path
// A local path usually means a path given as a relative or absolute filesystem path
func isLocalPath(p string) bool {
	re := regexp.MustCompile(`(?m)^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+`)
	return !re.MatchString(p)
}
