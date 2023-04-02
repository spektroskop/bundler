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
	"github.com/spektroskop/bundler/internal/tailwind"
)

type Bundler struct {
	Optimize    bool     `help:"Optimized build where applicable."`
	Tailwind    bool     `help:"Process stylesheets through tailwind"`
	Loaders     []string `help:"File loaders." placeholder:"EXTENSION"`
	Output      string   `help:"Output folder." placeholder:"PATH" required`
	Entrypoints []string `help:"Entrypoints to build." name:"entrypoint" arg`
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
	options.Plugins = []api.Plugin{
		elm.New(cli.Optimize),
		gleam.New(),
		gren.New(cli.Optimize),
		meta.New(),
	}

	if cli.Tailwind {
		options.Plugins = append(options.Plugins, tailwind.New())
	}

	options.Loader = make(map[string]api.Loader)
	for _, ext := range cli.Loaders {
		ext = fmt.Sprintf(".%s", ext)
		options.Loader[ext] = api.LoaderFile
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
