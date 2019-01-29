package gowatcher

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// NewWatcher new watcher
func NewWatcher(options WatcherOption) *Watcher {
	return &Watcher{
		Options: &options,
	}
}

//Start start
func (watcher *Watcher) Start() {

	if len(watcher.Options.Dirs) == 0 {
		//TODO: watcher service disabled
		return
	}

	fs, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

}
