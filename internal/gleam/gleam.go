package gleam

import (
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func New() api.Plugin {
	return api.Plugin{
		Name: "gleam",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.gleam$`},
				onResolve,
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "gleam"},
				onLoad(),
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

func onLoad() func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		var result api.OnLoadResult

		return result, nil
	}
}
