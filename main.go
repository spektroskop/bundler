package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s -output DIR [-optimize] [-debug] [-watch DIR] ENTRY...\n", os.Args[0])
		os.Exit(1)
	}

	output := flag.String("output", "", "")
	watch := flag.String("watch", "", "")
	files := flag.String("files", "", "")
	debug := flag.Bool("debug", false, "")
	optimize := flag.Bool("optimize", false, "")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	})

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if *output == "" || len(flag.Args()) == 0 {
		flag.Usage()
	}

	if err := os.MkdirAll(*output, 0750); err != nil {
		log.Fatal().Err(err).Msg("could not create output directory")
	}

	watcher, err := NewWatcher(*watch)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create watcher")
	}
	defer watcher.Close()

	var builds []*Build
	for _, path := range flag.Args() {
		build := NewBuild(path, *output, *optimize, strings.Split(*files, ","))
		builds = append(builds, build)
	}

	if *watch != "" {
		watcher.Watch(func(name string) {
			for _, build := range builds {
				build.Rebuild()
			}
		})
	}
}
