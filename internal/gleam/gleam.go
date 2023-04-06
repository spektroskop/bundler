package gleam

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func New(resolve string) api.Plugin {
	return api.Plugin{Name: "gleam", Setup: setup(resolve)}
}

func setup(resolve string) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `\.gleam$`}, onResolve)
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "gleam"}, onLoad(resolve))
	}
}

func onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	result := api.OnResolveResult{
		Namespace: "gleam", Path: filepath.Join(args.ResolveDir, args.Path),
	}

	return result, nil
}

func onLoad(resolve string) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		command, err := exec.LookPath("gleam")
		if err != nil {
			return api.OnLoadResult{}, err
		}

		cmd := exec.Command(command, "build", "--target=javascript")
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return api.OnLoadResult{}, err
		}

		var result api.OnLoadResult
		result.ResolveDir = resolve
		data, err := os.ReadFile(filepath.Join(result.ResolveDir,
			strings.Replace(filepath.Base(args.Path), ".gleam", ".mjs", -1),
		))
		if err == nil {
			contents := string(data)
			result.Contents = &contents
		}

		return result, err
	}
}
