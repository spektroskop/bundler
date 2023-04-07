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

func setup(path string) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnEnd(onEnd(path))
	}
}

func onEnd(path string) func(result *api.BuildResult) (api.OnEndResult, error) {
	return func(result *api.BuildResult) (r api.OnEndResult, _ error) {
		if len(result.Metafile) == 0 {
			return r, nil
		}

		data := []byte(result.Metafile)
		var meta Meta
		if err := json.Unmarshal(data, &meta); err != nil {
			return r, err
		}

		for name, output := range meta.Outputs {
			log.Info().Strs("source", source(output)).Msg(name)
		}

		return r, save(path, data)
	}
}

func source(output Output) (source []string) {
	for name := range output.Inputs {
		if parts := strings.Split(name, ":"); len(parts) == 2 {
			name = parts[1]
		}

		if wd, err := os.Getwd(); err == nil {
			if path, err := filepath.Rel(wd, name); err == nil {
				source = append(source, path)
				continue
			}
		}

		source = append(source, name)
	}

	return source
}

func save(path string, data []byte) error {
	if path == "" {
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	log.Info().Str("path", path).Msg("saved meta file")
	return nil
}
