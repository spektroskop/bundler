package main

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

func build(options api.BuildOptions) {
	result := api.Build(options)

	for _, msg := range result.Warnings {
		if msg.Location == nil {
			log.Warn().Msg(msg.Text)
		} else {
			log.Warn().Str("source", msg.Location.File).Msg(msg.Text)
		}
	}

	for _, msg := range result.Errors {
		if msg.Location == nil {
			log.Error().Msg(msg.Text)
		} else {
			log.Error().Str("source", msg.Location.File).Msg(msg.Text)
		}
	}

	if len(result.Errors) != 0 {
		os.Exit(1)
	}
}
