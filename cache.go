// Copyright (C) moneta. 2025-present.
//
// Created at 2025-08-02, by liasica

package godoc

import "os"

// GetGomodcache get gomodcache path from environment variables
func GetGomodcache() (gomodcache string) {
	gomodcache = os.Getenv("GOMODCACHE")
	if gomodcache == "" {
		gomodcache = os.Getenv("GOPATH")
		if gomodcache != "" {
			gomodcache = gomodcache + "/pkg/mod"
		}
	}
	return
}
