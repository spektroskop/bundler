package main

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/spektroskop/bundler/internal/elm"
	"github.com/spektroskop/bundler/internal/gleam"
	"github.com/spektroskop/bundler/internal/gren"
	"github.com/spektroskop/bundler/internal/meta"
	"github.com/spektroskop/bundler/internal/tailwind"
)

func configure(app App) (options api.BuildOptions) {
	options.Metafile = true
	options.EntryPoints = app.Entrypoints
	options.Bundle = true
	options.Outdir = app.Output
	options.EntryNames = "[dir]/[name]"
	options.AssetNames = "[dir]/[name]"
	options.Write = true
	options.MinifyWhitespace = app.Optimize
	options.MinifyIdentifiers = app.Optimize
	options.MinifySyntax = app.Optimize

	options.Loader = make(map[string]api.Loader)
	for ext, loader := range app.Loader {
		ext = fmt.Sprintf(".%s", ext)
		options.Loader[ext] = api.Loader(loader)
	}

	plugins := map[string]api.Plugin{
		"elm":   elm.New(app.Optimize),
		"gleam": gleam.New(app.Resolve["gleam"]),
		"gren":  gren.New(app.Optimize),
	}

	for _, name := range app.Activate {
		switch name {
		case "tailwind":
			plugins[name] = tailwind.New()
		}
	}

	for _, name := range app.Deactivate {
		delete(plugins, name)
	}

	options.Plugins = []api.Plugin{meta.New(app.Meta)}
	for _, plugin := range plugins {
		options.Plugins = append(options.Plugins, plugin)
	}

	return options
}
