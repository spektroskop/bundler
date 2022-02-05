package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debounce := flag.Duration("debounce", 250*time.Millisecond, "")
	debug := flag.Bool("debug", false, "")
	files := flag.String("files", "", "")
	notifyMethod := flag.String("notify-method", "PATCH", "")
	notifyURL := flag.String("notify-url", "", "")
	optimize := flag.Bool("optimize", false, "")
	output := flag.String("output", "", "")
	watch := flag.String("watch", "", "")

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
		maybeNotify(*notifyMethod, *notifyURL)

		watcher.Watch(*debounce, *notifyMethod, *notifyURL, func(name string) {
			for _, build := range builds {
				build.Rebuild()
			}
		})
	}
}
