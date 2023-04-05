package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spektroskop/bundler/internal/elm"
	"github.com/spektroskop/bundler/internal/gleam"
	"github.com/spektroskop/bundler/internal/gren"
	"github.com/spektroskop/bundler/internal/meta"
	"github.com/spektroskop/bundler/internal/plugin"
	"github.com/spektroskop/bundler/internal/tailwind"
)

type Bundler struct {
	Activate    []string `help:"List of optional plugins to activate (${enum})." enum:"tailwind" placeholder:"NAME"`
	Deactivate  []string `help:"List of plugins to deactivate (${enum})." enum:"elm,gleam,gren" placeholder:"NAME"`
	Entrypoints []string `help:"Entrypoints to build." name:"entrypoint" arg`
	Loaders     []string `help:"File loaders." placeholder:"EXTENSION"`
	Meta        string   `help:"Meta file output." placeholder:"PATH"`
	Optimize    bool     `help:"Optimized build where applicable."`
	Output      string   `help:"Output folder." placeholder:"PATH" required`
	Resolve     string   `help:"Import resolve dir" placeholder:"PATH"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stdout, PartsExclude: []string{"time"},
	})

	var cli Bundler
	kong.Parse(&cli)

	var options api.BuildOptions

	options.Metafile = true
	options.EntryPoints = cli.Entrypoints
	options.Bundle = true
	options.Outdir = cli.Output
	options.EntryNames = "[dir]/[name]"
	options.AssetNames = "[dir]/[name]"
	options.Write = true
	options.MinifyWhitespace = cli.Optimize
	options.MinifyIdentifiers = cli.Optimize
	options.MinifySyntax = cli.Optimize

	options.Loader = make(map[string]api.Loader)
	for _, ext := range cli.Loaders {
		ext = fmt.Sprintf(".%s", ext)
		options.Loader[ext] = api.LoaderFile
	}

	options.Plugins = []api.Plugin{meta.New(cli.Meta)}
	config := plugin.Config{Optimized: cli.Optimize, Resolve: cli.Resolve}

	plugins := map[string]api.Plugin{
		"elm":   elm.New(config),
		"gleam": gleam.New(config),
		"gren":  gren.New(config),
	}

	for _, name := range cli.Activate {
		switch name {
		case "tailwind":
			plugins[name] = tailwind.New(config)
		}
	}

	for _, name := range cli.Deactivate {
		delete(plugins, name)
	}

	for _, plugin := range plugins {
		options.Plugins = append(options.Plugins, plugin)
	}

	result := api.Build(options)
	formatOptions := api.FormatMessagesOptions{Color: true}

	for _, msg := range api.FormatMessages(result.Warnings, formatOptions) {
		fmt.Print(msg)
	}

	for _, msg := range api.FormatMessages(result.Errors, formatOptions) {
		fmt.Print(msg)
	}

	if len(result.Errors) != 0 {
		os.Exit(1)
	}
}
