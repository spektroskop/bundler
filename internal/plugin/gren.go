package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func Gren(config Config) api.Plugin {
	var plugin api.Plugin
	plugin.Name = "gren"

	plugin.Setup = func(build api.PluginBuild) {
		var resolveOptions api.OnResolveOptions
		resolveOptions.Filter = `\.gren$`

		build.OnResolve(
			resolveOptions,
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				var result api.OnResolveResult
				result.Path = filepath.Join(args.ResolveDir, args.Path)
				result.Namespace = "gren"
				return result, nil
			},
		)

		var loadOptions api.OnLoadOptions
		loadOptions.Filter = `.*`
		loadOptions.Namespace = "gren"
		build.OnLoad(loadOptions, gren(config))
	}

	return plugin
}

func gren(config Config) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult
		result.ResolveDir = config.Resolve

		command, err := exec.LookPath("gren")
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
