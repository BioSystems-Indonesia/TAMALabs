package web

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed dist/*
var content embed.FS

// Content returns the embedded content.
func Content() fs.FS {
	dir, err := fs.Sub(content, "dist")
	if err != nil {
		panic(fmt.Errorf("could not find dist directory: %w", err))
	}
	return dir
}
