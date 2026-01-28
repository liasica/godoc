// Copyright (C) godoc. 2026-present.
//
// Created at 2026-01-09, by liasica

package godoc

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	basePath string

	MainFile         string            `yaml:"mainFile"`
	MarkdownFilesDir string            `yaml:"markdownFilesDir"`
	ExternalPath     map[string]string `yaml:"externalPath"`
	Path             []string          `yaml:"path"`
	Output           string            `yaml:"output"`
	OutputTypes      []string          `yaml:"outputTypes"`
}

func defaultConfig() *Config {
	return &Config{
		ExternalPath: map[string]string{},
		Path: []string{
			"./internal/interface/http/router",
			"./internal/interface/http/handler",
			"./internal/infrastructure/model",
			"./internal/infrastructure/dao/pagination",
			"./internal/presentation/dto",
		},
		MainFile: "route.go",
		Output:   "./assets/docs/",
		OutputTypes: []string{
			"yaml",
		},
	}
}

// DefaultConfig returns the default configuration as a YAML string.
func DefaultConfig() string {
	defaultCfg := defaultConfig()
	b, _ := yaml.Marshal(defaultCfg)
	return string(b)
}

// LoadConfig reads a YAML config file. If path is empty, uses ".godoc.yaml" in cwd.
// It returns a Config with defaults applied when fields are missing.
func LoadConfig(path string) (cfg *Config, err error) {
	if path == "" {
		path = ".godoc.yaml"
	}

	var b []byte
	b, err = os.ReadFile(path)
	if err != nil {
		// If file not found, return defaults rather than error
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, err
	}

	cfg = defaultConfig()
	if err = yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	cfg.basePath = filepath.Dir(path)

	return
}

// ResolveConfigPath returns an absolute path for config if provided relative.
func ResolveConfigPath(p string) string {
	if p == "" {
		p = ".godoc.yaml"
	}
	if !filepath.IsAbs(p) {
		abs, err := filepath.Abs(p)
		if err == nil {
			p = abs
		}
	}
	return p
}

// ResolveAbsPath resolves dir relative to the config file's base path.
func (c *Config) ResolveAbsPath(dir string) (string, error) {
	return filepath.Abs(filepath.Join(c.basePath, dir))
}
