package elm

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func New(optimize bool) api.Plugin {
	return api.Plugin{Name: "elm", Setup: setup(optimize)}
}

func setup(optimize bool) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `\.elm$`}, onResolve)
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "elm"}, onLoad(optimize))
	}
}

func onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	result := api.OnResolveResult{
		Namespace: "elm", Path: filepath.Join(args.ResolveDir, args.Path),
	}

	return result, nil
}

func onLoad(optimize bool) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		command, err := exec.LookPath("elm")
		if err != nil {
			return api.OnLoadResult{}, err
		}

		output, err := os.CreateTemp("/tmp", "*.js")
		if err != nil {
			return api.OnLoadResult{}, err
		}
		defer os.Remove(output.Name())

		parts := []string{command, "make", args.Path, "--output", output.Name()}

		if optimize {
			parts = append(parts, "--optimize")
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return api.OnLoadResult{}, err
		}

		var result api.OnLoadResult
		data, err := os.ReadFile(output.Name())
		if err == nil {
			contents := string(data)
			result.Contents = &contents
		}

		return result, err
	}
}
