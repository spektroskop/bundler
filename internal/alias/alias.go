package alias

import (
	"github.com/evanw/esbuild/pkg/api"
)

const (
	Filter    = `.*`
	Name      = "alias"
	Namespace = "alias"
)

func New(aliases map[string]string) api.Plugin {
	return api.Plugin{Name: Name, Setup: setup(aliases)}
}

func setup(aliases map[string]string) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(
			api.OnResolveOptions{Filter: Filter},
			func(args api.OnResolveArgs) (r api.OnResolveResult, _ error) {
				if alias, ok := aliases[args.Path]; ok {
					r.Namespace = "fetch"
					r.Path = alias
				}

				return r, nil
			},
		)
	}
}
