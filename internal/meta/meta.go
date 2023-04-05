package meta

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

type Output struct {
	Inputs map[string]any `json:"inputs"`
}

type Meta struct {
	Outputs map[string]Output `json:"outputs"`
}

func New(save string) api.Plugin {
	return api.Plugin{Name: "meta", Setup: setup(save)}
}

func setup(save string) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnEnd(onEnd(save))
	}
}

func onEnd(save string) func(result *api.BuildResult) (api.OnEndResult, error) {
	return func(result *api.BuildResult) (api.OnEndResult, error) {
		if len(result.Metafile) == 0 {
			return api.OnEndResult{}, nil
		}

		var meta Meta
		if err := json.Unmarshal([]byte(result.Metafile), &meta); err != nil {
			return api.OnEndResult{}, err
		}

		if save != "" {
			f, err := os.Create(save)
			if err != nil {
				return api.OnEndResult{}, err
			}

			if _, err := f.Write([]byte(result.Metafile)); err != nil {
				return api.OnEndResult{}, err
			}

			log.Info().Str("path", save).Msg("saved meta file")
		}

		for name, output := range meta.Outputs {
			var source []string

			for name := range output.Inputs {
				if parts := strings.Split(name, ":"); len(parts) == 2 {
					if wd, err := os.Getwd(); err == nil {
						if path, err := filepath.Rel(wd, parts[1]); err == nil {
							source = append(source, path)
							continue
						}
					}
				}

				source = append(source, name)
			}

			log.Info().Strs("source", source).Msg(name)
		}

		return api.OnEndResult{}, nil
	}
}
