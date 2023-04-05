package gleam

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func New(path string) api.Plugin {
	return api.Plugin{
		Name: "gleam",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.gleam$`},
				onResolve,
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "gleam"},
				onLoad(path),
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

func onLoad(path string) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult
		result.ResolveDir = path

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

		compiled, err := os.ReadFile(entrypoint)
		if err == nil {
			contents := string(compiled)
			result.Contents = &contents
		}

		return result, err
	}
}
