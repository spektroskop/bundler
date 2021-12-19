package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

func NewElmPlugin(optimize bool) api.Plugin {
	return api.Plugin{
		Name: "elm",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.elm$`},
				ElmOnResolve,
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "elm"},
				ElmOnLoad(optimize),
			)
		},
	}
}

func ElmOnResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	result := api.OnResolveResult{
		Path:      filepath.Join(args.ResolveDir, args.Path),
		Namespace: "elm",
	}

	return result, nil
}

func ElmOnLoad(optimize bool) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult

		if _, err := exec.LookPath("elm"); err != nil {
			return result, err
		}

		temp, err := os.CreateTemp("/tmp", "*.js")
		if err != nil {
			return result, err
		}
		defer os.Remove(temp.Name())

		wd, err := os.Getwd()
		if err != nil {
			return result, err
		}

		path, err := filepath.Rel(wd, args.Path)
		if err != nil {
			return result, err
		}

		elmMake := []string{"elm", "make"}

		if optimize {
			elmMake = append(elmMake, "--optimize")
		}

		elmMake = append(elmMake, path, fmt.Sprintf("--output=%s", temp.Name()))
		cmd := exec.Command(elmMake[0], elmMake[1:]...)
		cmd.Stderr = os.Stderr

		log.Debug().Str("path", path).Msg("compile")

		if err := cmd.Run(); err != nil {
			return result, err
		}

		compiled, err := os.ReadFile(temp.Name())
		if err != nil {
			return result, err
		}

		contents := string(compiled)
		result.Contents = &contents
		return result, nil
	}
}
