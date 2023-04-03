package tailwind

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func New() api.Plugin {
	return api.Plugin{
		Name: "tailwind",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.css$`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					path := filepath.Join(args.ResolveDir, args.Path)
					return api.OnResolveResult{Path: path, Namespace: "tailwind"}, nil
				},
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "tailwind"},
				onLoad(),
			)
		},
	}
}

func onLoad() func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult
		result.Loader = api.LoaderCSS

		command, err := exec.LookPath("tailwind")
		if err != nil {
			return result, err
		}

		wd, err := os.Getwd()
		if err != nil {
			return result, err
		}

		path, err := filepath.Rel(wd, args.Path)
		if err != nil {
			return result, err
		}

		parts := []string{command, "--input", path}

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
