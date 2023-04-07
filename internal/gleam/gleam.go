package gleam

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func New(resolve string) api.Plugin {
	return api.Plugin{Name: "gleam", Setup: setup(resolve)}
}

func setup(resolve string) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `\.gleam$`}, onResolve)
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "gleam"}, onLoad(resolve))
	}
}

func onResolve(args api.OnResolveArgs) (r api.OnResolveResult, _ error) {
	r.Namespace = "gleam"
	r.Path = filepath.Join(args.ResolveDir, args.Path)
	return r, nil
}

func onLoad(resolve string) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (r api.OnLoadResult, _ error) {
		command, err := exec.LookPath("gleam")
		if err != nil {
			return r, err
		}

		cmd := exec.Command(command, "build", "--target=javascript")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return r, err
		}

		r.ResolveDir = resolve
		data, err := os.ReadFile(filepath.Join(r.ResolveDir,
			strings.Replace(filepath.Base(args.Path), ".gleam", ".mjs", -1),
		))
		if err == nil {
			contents := string(data)
			r.Contents = &contents
		}

		return r, err
	}
}
