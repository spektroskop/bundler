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
	return api.Plugin{Name: "tailwind", Setup: setup()}
}

func setup() func(build api.PluginBuild) {
	return func(build api.PluginBuild) {
		build.OnResolve(api.OnResolveOptions{Filter: `\.css$`}, onResolve)
		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "tailwind"}, onLoad())
	}
}

func onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	result := api.OnResolveResult{
		Path: filepath.Join(args.ResolveDir, args.Path), Namespace: "tailwind",
	}

	return result, nil
}

func onLoad() func(api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		command, err := exec.LookPath("tailwind")
		if err != nil {
			return api.OnLoadResult{}, err
		}

		var stderr bytes.Buffer
		cmd := exec.Command(command, "--input", args.Path)
		cmd.Stderr = &stderr

		data, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, stderr.String())
			return api.OnLoadResult{}, err
		}

		contents := string(data)
		var result api.OnLoadResult
		result.Loader = api.LoaderCSS
		result.Contents = &contents

		return result, nil
	}
}
