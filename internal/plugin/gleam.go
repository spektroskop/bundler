package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func Gleam(config Config) api.Plugin {
	var plugin api.Plugin
	plugin.Name = "gleam"

	plugin.Setup = func(build api.PluginBuild) {
		var resolveOptions api.OnResolveOptions
		resolveOptions.Filter = `\.gleam$`

		build.OnResolve(
			resolveOptions,
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				var result api.OnResolveResult
				result.Path = filepath.Join(args.ResolveDir, args.Path)
				result.Namespace = "gleam"
				return result, nil
			},
		)

		var loadOptions api.OnLoadOptions
		loadOptions.Filter = `.*`
		loadOptions.Namespace = "gleam"
		build.OnLoad(loadOptions, gleam(config))
	}

	return plugin
}

func gleam(config Config) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult
		result.ResolveDir = config.Resolve

		// FIXME: Hm..
		source := filepath.Base(args.Path)
		entrypoint := fmt.Sprintf("%s/%s",
			result.ResolveDir,
			strings.Replace(source, "gleam", "mjs", -1),
		)

		command, err := exec.LookPath("gleam")
		if err != nil {
			return result, err
		}

		parts := []string{
			command, "build", "--target=javascript",
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return result, err
		}

		compiled, err := os.ReadFile(entrypoint)
		if err == nil {
			contents := string(compiled)
			result.Contents = &contents
		}

		return result, err
	}
}
