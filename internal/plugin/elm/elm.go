package elm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

func New(optimize bool) api.Plugin {
	return api.Plugin{
		Name: "elm",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.elm$`},
				onResolve,
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "elm"},
				onLoad(optimize),
			)
		},
	}
}

func onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	result := api.OnResolveResult{
		Path:      filepath.Join(args.ResolveDir, args.Path),
		Namespace: "elm",
	}

	return result, nil
}

func onLoad(optimize bool) func(api.OnLoadArgs) (api.OnLoadResult, error) {
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

		buildCommand := []string{"elm", "make"}

		if optimize {
			buildCommand = append(buildCommand, "--optimize")
		}

		buildCommand = append(buildCommand, path, fmt.Sprintf("--output=%s", temp.Name()))
		cmd := exec.Command(buildCommand[0], buildCommand[1:]...)
		cmd.Stderr = os.Stderr

		log.Info().Str("plugin", "elm").Str("path", path).Msg("build")

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
