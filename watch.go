package main

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/rs/zerolog/log"
	"gopkg.in/fsnotify.v1"
)

type Watcher struct {
	*fsnotify.Watcher
}

func NewWatcher(root string) (*Watcher, error) {
	var err error
	watcher := new(Watcher)

	if watcher.Watcher, err = fsnotify.NewWatcher(); err != nil {
		return watcher, err
	}

	if strings.TrimSpace(root) == "" {
		return watcher, nil
	}

	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			return nil
		} else if strings.HasPrefix(entry.Name(), ".") {
			return filepath.SkipDir
		}

		log.Debug().Str("path", path).Msg("watch folder")
		if err := watcher.Add(path); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return watcher, err
	}

	return watcher, nil
}

func (watcher *Watcher) Watch(onChange func(string)) {
	debounced := debounce.New(250 * time.Millisecond)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				debounced(func() {
					log.Info().Str("path", event.Name).Msg("changed")
					onChange(event.Name)
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Error().Err(err).Msg("watch error")
		}
	}
}
