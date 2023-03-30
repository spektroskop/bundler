package plugin

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

func Analyze() api.Plugin {
	onEnd := func(result *api.BuildResult) (api.OnEndResult, error) {
		if len(result.Metafile) == 0 {
			return api.OnEndResult{}, nil
		}

		var meta Meta
		if err := json.Unmarshal([]byte(result.Metafile), &meta); err != nil {
			return api.OnEndResult{}, err
		}

		for name, output := range meta.Outputs {
			var inputs []string
			for name := range output.Inputs {
				if strings.Contains(name, ":") {
					if parts := strings.Split(name, ":"); len(parts) == 2 {
						if wd, err := os.Getwd(); err == nil {
							if rel, err := filepath.Rel(wd, parts[1]); err == nil {
								inputs = append(inputs, rel)
								continue
							}
						}
					}
				}

				inputs = append(inputs, name)
			}

			log.Info().
				Strs("source", inputs).
				Msg(name)
		}

		return api.OnEndResult{}, nil
	}

	return api.Plugin{
		Name: "PostBuildActions",
		Setup: func(build api.PluginBuild) {
			build.OnEnd(onEnd)
		},
	}
}