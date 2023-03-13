package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Bundler struct {
	Optimize   bool     `help:"Optimized build where applicable."`
	Debug      bool     `help:"Show debug messages."`
	Load       []string `help:"File loaders." placeholder:"ext"`
	Output     string   `help:"Output folder." required`
	Entrypoint []string `help:"Entrypoints to build." arg placeholder:"path"`
}

func main() {
	var cli Bundler
	kong.Parse(&cli)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stdout, PartsExclude: []string{"time"},
	})

	if cli.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().
		Bool("debug", cli.Debug).
		Strs("load", cli.Load).
		Bool("optimize", cli.Optimize).
		Str("output", cli.Output).
		Strs("entrypoints", cli.Entrypoint).
		Msg("bundle")

	if err := os.MkdirAll(cli.Output, 0750); err != nil {
		log.Fatal().Err(err).Msg("could not create output directory")
	}

	var builds []*Build
	for _, path := range cli.Entrypoint {
		build := NewBuild(path, cli.Output, cli.Optimize, cli.Load)
		builds = append(builds, build)
	}
}
