package tailwind

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/spektroskop/bundler/internal/plugin"
)

func New(config plugin.Config) api.Plugin {
	return api.Plugin{Name: "tailwind", Setup: setup(config)}
}

func setup(config plugin.Config) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(
			api.OnResolveOptions{Filter: `\.css$`},
			onResolve,
		)

		build.OnLoad(
			api.OnLoadOptions{Filter: `.*`, Namespace: "tailwind"},
			onLoad(config),
		)
	}
}

func onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	var result api.OnResolveResult
	result.Path = filepath.Join(args.ResolveDir, args.Path)
	result.Namespace = "tailwind"
	return result, nil
}

func onLoad(config plugin.Config) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult
		result.ResolveDir = config.Resolve
		result.Loader = api.LoaderCSS

		command, err := exec.LookPath("tailwind")
		if err != nil {
			return result, err
		}

		parts := []string{command, "--input", args.Path}

		var stderr bytes.Buffer
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = &stderr

		compiled, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, stderr.String())
			return result, err
		}

		contents := string(compiled)
		result.Contents = &contents

		return result, nil
	}
}
