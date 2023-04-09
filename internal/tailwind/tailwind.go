package tailwind

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func New(config string) api.Plugin {
	return api.Plugin{Name: "tailwind", Setup: setup(config)}
}

func setup(config string) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `\.css$`}, onResolve)
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "tailwind"}, onLoad(config))
	}
}

func onResolve(args api.OnResolveArgs) (r api.OnResolveResult, _ error) {
	r.Path = filepath.Join(args.ResolveDir, args.Path)
	r.Namespace = "tailwind"
	return r, nil
}

func onLoad(config string) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (r api.OnLoadResult, _ error) {
		command, err := exec.LookPath("tailwind")
		if err != nil {
			return r, err
		}

		parts := []string{command, "--input", args.Path}
		if config != "" {
			parts = append(parts, "--config", config)
		}

		var stderr bytes.Buffer
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = &stderr

		data, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, stderr.String())
			return r, err
		}

		contents := string(data)
		r.Contents = &contents
		r.Loader = api.LoaderCSS
		return r, nil
	}
}
