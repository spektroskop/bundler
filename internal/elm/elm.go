package elm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/spektroskop/bundler/internal/plugin"
)

func New(config plugin.Config) api.Plugin {
	return api.Plugin{Name: "elm", Setup: setup(config)}
}

func setup(config plugin.Config) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(
			api.OnResolveOptions{Filter: `\.elm$`},
			onResolve,
		)

		build.OnLoad(
			api.OnLoadOptions{Filter: `.*`, Namespace: "elm"},
			onLoad(config),
		)
	}
}

func onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	var result api.OnResolveResult
	result.Path = filepath.Join(args.ResolveDir, args.Path)
	result.Namespace = "elm"
	return result, nil
}

func onLoad(config plugin.Config) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult
		result.ResolveDir = config.Resolve

		command, err := exec.LookPath("elm")
		if err != nil {
			return result, err
		}

		output, err := os.CreateTemp("/tmp", "*.js")
		if err != nil {
			return result, err
		}
		defer os.Remove(output.Name())

		parts := []string{
			command, "make", args.Path,
			fmt.Sprintf("--output=%s", output.Name()),
		}

		if config.Optimized {
			parts = append(parts, "--optimize")
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return result, err
		}

		compiled, err := os.ReadFile(output.Name())
		if err == nil {
			contents := string(compiled)
			result.Contents = &contents
		}

		return result, err
	}
}
