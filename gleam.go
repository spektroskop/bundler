package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/rs/zerolog/log"
)

func NewGleamPlugin() api.Plugin {
	return api.Plugin{
		Name: "gleam",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.gleam$`},
				GleamOnResolve,
			)

			build.OnLoad(
				api.OnLoadOptions{Filter: `.*`, Namespace: "gleam"},
				GleamOnLoad,
			)
		},
	}
}

func GleamOnResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	result := api.OnResolveResult{
		Path:      filepath.Join(args.ResolveDir, args.Path),
		Namespace: "gleam",
	}

	return result, nil
}

func GleamOnLoad(args api.OnLoadArgs) (api.OnLoadResult, error) {
	var result api.OnLoadResult

	if _, err := exec.LookPath("gleam"); err != nil {
		return result, err
	}

	temp, err := os.MkdirTemp("/tmp", "gleam")
	if err != nil {
		return result, err
	}
	defer os.Remove(temp)

	wd, err := os.Getwd()
	if err != nil {
		return result, err
	}

	path, err := filepath.Rel(wd, args.Path)
	if err != nil {
		return result, err
	}

	gleamCmd := []string{
		"gleam",
		"compile-package",
		"--name",
		"TODO",
		"--target",
		"javascript",
		"--src",
		path,
		"--out",
		temp,
	}

	cmd := exec.Command(gleamCmd[0], gleamCmd[1:]...)
	cmd.Stderr = os.Stderr

	log.Debug().Str("path", path).Msg("compile")

	if err := cmd.Run(); err != nil {
		return result, err
	}

	contents := ""
	result.Contents = &contents

	return result, nil
}
