package gleam

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func New(resolveDir string) api.Plugin {
	return api.Plugin{
		Name: "gleam",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.gleam$`},
				onResolve,
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "gleam"},
				onLoad(resolveDir),
			)
		},
	}
}

func onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	result := api.OnResolveResult{
		Path:      filepath.Join(args.ResolveDir, args.Path),
		Namespace: "gleam",
	}

	return result, nil
}

func onLoad(resolveDir string) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult
		result.ResolveDir = resolveDir

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
