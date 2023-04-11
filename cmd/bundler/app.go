package main

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
)

type App struct {
	Entrypoints []string          `help:"Entrypoints to build." name:"entrypoint" arg`
	Output      string            `help:"Output folder." placeholder:"PATH" required`
	Optimized   bool              `help:"Optimized build where applicable."`
	Meta        string            `help:"Meta file output." placeholder:"PATH"`
	Activate    []string          `help:"List of optional plugins to activate (${enum})." enum:"tailwind" placeholder:"NAME"`
	Deactivate  []string          `help:"List of plugins to deactivate (${enum})." enum:"elm,gleam,gren" placeholder:"NAME"`
	Loader      map[string]Loader `help:"Loaders (jsx,file)." placeholder:"EXT:NAME"`
	Config      map[string]string `help:"Set config values." name:"set" placeholder:"KEY=VALUE" mapsep:","` // TODO: List config keys in help text
}

type Loader api.Loader

func (v *Loader) UnmarshalText(b []byte) error {
	switch string(b) {
	case "jsx":
		*v = Loader(api.LoaderJSX)
		return nil
	case "file":
		*v = Loader(api.LoaderFile)
		return nil
	default:
		return fmt.Errorf("bad loader: %s", string(b))
	}
}
