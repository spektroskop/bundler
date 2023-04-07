package elm

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func New(optimize bool) api.Plugin {
	return api.Plugin{Name: "elm", Setup: setup(optimize)}
}

func setup(optimize bool) func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `\.elm$`}, onResolve)
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "elm"}, onLoad(optimize))
	}
}

func onResolve(args api.OnResolveArgs) (r api.OnResolveResult, _ error) {
	r.Namespace = "elm"
	r.Path = filepath.Join(args.ResolveDir, args.Path)
	return r, nil
}

func onLoad(optimize bool) func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (r api.OnLoadResult, _ error) {
		command, err := exec.LookPath("elm")
		if err != nil {
			return r, err
		}

		output, err := os.CreateTemp("/tmp", "*.js")
		if err != nil {
			return r, err
		}
		defer os.Remove(output.Name())

		parts := []string{command, "make", args.Path, "--output", output.Name()}
		if optimize {
			parts = append(parts, "--optimize")
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return r, err
		}

		data, err := os.ReadFile(output.Name())
		if err == nil {
			contents := string(data)
			r.Contents = &contents
		}

		return r, err
	}
}
