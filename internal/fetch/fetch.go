package fetch

import (
	"io"
	"net/http"
	"net/url"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

const Filter = `^https?://`
const Namespace = "fetch"

func New() api.Plugin {
	return api.Plugin{Name: "fetch", Setup: setup()}
}

func setup() func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: Filter}, func(args api.OnResolveArgs) (r api.OnResolveResult, _ error) {
			r.Namespace = Namespace
			r.Path = args.Path
			return r, nil
		})

		build.OnResolve(api.OnResolveOptions{Filter: ".*", Namespace: Namespace}, func(args api.OnResolveArgs) (r api.OnResolveResult, _ error) {
			base, err := url.Parse(args.Importer)
			if err != nil {
				return r, err
			}

			relative, err := url.Parse(args.Path)
			if err != nil {
				return r, err
			}

			r.Namespace = Namespace
			r.Path = base.ResolveReference(relative).String() // FIXME

			log.Info().Str("path", r.Path).Msg("resolve")
			return r, nil
		})

		build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: Namespace}, func(args api.OnLoadArgs) (r api.OnLoadResult, _ error) {
			log.Info().Str("path", args.Path).Msg("load")

			res, err := http.Get(args.Path)
			if err != nil {
				return r, err
			}
			defer res.Body.Close()

			bytes, err := io.ReadAll(res.Body)
			if err != nil {
				return r, err
			}

			contents := string(bytes)
			r.Contents = &contents
			return r, nil
		})
	}
}
