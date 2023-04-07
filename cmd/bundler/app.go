package main

type App struct {
	Entrypoints []string          `help:"Entrypoints to build." name:"entrypoint" arg`
	Output      string            `help:"Output folder." placeholder:"PATH" required`
	Optimize    bool              `help:"Optimized build where applicable."`
	Meta        string            `help:"Meta file output." placeholder:"PATH"`
	Activate    []string          `help:"List of optional plugins to activate (${enum})." enum:"tailwind" placeholder:"NAME"`
	Deactivate  []string          `help:"List of plugins to deactivate (${enum})." enum:"elm,gleam,gren" placeholder:"NAME"`
	Loaders     []string          `help:"File loaders." placeholder:"EXTENSION"`
	Resolve     map[string]string `help:"Plugin resolve path." placeholder:"PLUGIN=PATH"`
}
