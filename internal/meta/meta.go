package meta

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

type Input struct{}

type Output struct {
	Inputs map[string]Input `json:"inputs"`
}

type Meta struct {
	Outputs map[string]Output `json:"outputs"`
}

func New(path string) api.Plugin {
	return api.Plugin{
		Name: "meta",
		Setup: func(build api.PluginBuild) {
			build.OnEnd(onEnd(path))
		},
	}
}

func getInputs(output Output) []string {
	var names []string

	for name := range output.Inputs {
		if parts := strings.Split(name, ":"); len(parts) == 2 {
			if wd, err := os.Getwd(); err == nil {
				if path, err := filepath.Rel(wd, parts[1]); err == nil {
					names = append(names, path)
					continue
				}
			}
		}

		names = append(names, name)
	}

	return names
}

func saveMeta(path string, data string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(data)); err != nil {
		return err
	}

	log.Info().Str("path", path).Msg("saved meta file")
	return nil
}

func onEnd(path string) func(result *api.BuildResult) (api.OnEndResult, error) {
	return func(result *api.BuildResult) (api.OnEndResult, error) {
		if len(result.Metafile) == 0 {
			return api.OnEndResult{}, nil
		}

		var meta Meta
		if err := json.Unmarshal([]byte(result.Metafile), &meta); err != nil {
			return api.OnEndResult{}, err
		}

		if path != "" {
			if err := saveMeta(path, result.Metafile); err != nil {
				return api.OnEndResult{}, err
			}
		}

		for name, output := range meta.Outputs {
			log.Info().Strs("source", getInputs(output)).Msg(name)
		}

		return api.OnEndResult{}, nil
	}
}
