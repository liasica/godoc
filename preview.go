// Copyright (C) godoc. 2026-present.
//
// Created at 2026-01-12, by liasica

package godoc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/liasica/godoc/assets"
)

func GetAvailableAddress() (string, error) {
	// Listen on localhost with port 0 to get an available port from the OS.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	defer func() { _ = l.Close() }()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return "", fmt.Errorf("failed to get tcp address from listener")
	}

	return fmt.Sprintf("127.0.0.1:%d", addr.Port), nil
}

// Preview Starts a local server to preview the generated documentation.
func Preview(cfgPath, address string) error {
	basePath := filepath.Dir(cfgPath)

	// parse config
	cfg, err := LoadConfig(ResolveConfigPath(cfgPath))
	if err != nil {
		return err
	}

	e := echo.New()
	e.Renderer = LoadTemplates(assets.TemplateFS, "templates")
	e.HideBanner = true

	e.GET("/docs/openapi.yaml", func(c echo.Context) error {
		f := filepath.Join(basePath, cfg.Output, "swagger.yaml")
		return c.File(f)
	})

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "docs.html", map[string]string{
			"title":   "Go Documentation Preview",
			"specURL": "/docs/openapi.yaml",
		})
	})

	// Gracefully start the server
	go func() {
		fmt.Printf("starting preview server at http://%s\n", address)
		err = e.Start(address)
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			fmt.Printf("failed to start preview server at http://%s\n", address)
			os.Exit(1)
		}
	}()

	// Open the URL in the default browser
	err = OpenURL(fmt.Sprintf("http://%s", address))
	if err != nil {
		fmt.Printf("failed to open browser: %v\n", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = e.Shutdown(ctx)
	if err != nil {
		fmt.Printf("error shutting down HTTP server: %v\n", err)
		os.Exit(1)
	}

	return nil
}
