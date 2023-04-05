package plugin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func Tailwind(config Config) api.Plugin {
	var plugin api.Plugin
	plugin.Name = "tailwind"

	plugin.Setup = func(build api.PluginBuild) {
		var resolveOptions api.OnResolveOptions
		resolveOptions.Filter = `\.css$`

		build.OnResolve(
			resolveOptions,
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				var result api.OnResolveResult
				result.Path = filepath.Join(args.ResolveDir, args.Path)
				result.Namespace = "css"
				return result, nil
			},
		)

		var loadOptions api.OnLoadOptions
		loadOptions.Filter = `.*`
		loadOptions.Namespace = "css"
		build.OnLoad(loadOptions, tailwind(config))
	}

	return plugin
}

func tailwind(config Config) func(api.OnLoadArgs) (api.OnLoadResult, error) {
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
