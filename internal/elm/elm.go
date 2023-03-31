package elm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func New(optimize bool) api.Plugin {
	return api.Plugin{
		Name: "elm",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.elm$`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					path := filepath.Join(args.ResolveDir, args.Path)
					return api.OnResolveResult{Path: path, Namespace: "elm"}, nil
				},
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "elm"},
				onLoad(optimize),
			)
		},
	}
}

func onLoad(optimize bool) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult

		if _, err := exec.LookPath("elm"); err != nil {
			return result, err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return result, err
		}

		path, err := filepath.Rel(cwd, args.Path)
		if err != nil {
			return result, err
		}

		output, err := os.CreateTemp("/tmp", "*.js")
		if err != nil {
			return result, err
		}
		defer os.Remove(output.Name())

		parts := []string{"elm", "make", path, fmt.Sprintf("--output=%s", output.Name())}
		if optimize {
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
