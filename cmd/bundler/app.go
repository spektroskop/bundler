package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/evanw/esbuild/pkg/api"
)

const (
	ConfigGleamResolve   = "gleam.resolve"
	ConfigMetaOutput     = "meta.output"
	ConfigTailwindConfig = "tailwind.config"
)

var (
	configs = map[string]string{
		fmt.Sprintf("%s=PATH", ConfigGleamResolve):   "Path to use for resolving dependencies in Gleam.",
		fmt.Sprintf("%s=PATH", ConfigMetaOutput):     "Save build metadata to file.",
		fmt.Sprintf("%s=PATH", ConfigTailwindConfig): "Path to custom tailwind config file.",
	}
)

type App struct {
	Entrypoints []string          `help:"Entrypoints to build." name:"entrypoint" arg`
	Output      string            `help:"Output folder." short:"o" placeholder:"PATH" required`
	Optimized   bool              `help:"Optimized build where applicable." short:"z"`
	Activate    []string          `help:"List of optional plugins to activate (${enum})." short:"a" enum:"tailwind" placeholder:"NAME"`
	Deactivate  []string          `help:"List of plugins to deactivate (${enum})." short:"d" enum:"elm,gleam,gren" placeholder:"NAME"`
	Loader      map[string]Loader `help:"Loaders (jsx,file)." short:"l" placeholder:"EXT:NAME"`
	Config      map[string]string `help:"Set config values." short:"s" name:"set" placeholder:"KEY=VALUE" mapsep:","`
}

func (app App) Help(options kong.HelpOptions, ctx *kong.Context) error {
	if err := kong.DefaultHelpPrinter(options, ctx); err != nil {
		return err
	}

	var pad int
	for key := range configs {
		if len(key) > pad {
			pad = len(key)
		}
	}

	fmt.Println()
	fmt.Println("Config:")
	for key, help := range configs {
		fmt.Printf("  %-*s    %s\n", pad, key, help)
	}
	fmt.Println()

	return nil
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
