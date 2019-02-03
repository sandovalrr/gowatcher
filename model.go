package gowatcher

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

//Repo repo
var Repo = "github.com/sandovalrr/gowatcher"

// WatcherOption Options Model
type WatcherOption struct {
	Dirs       []string
	Recursive  bool
	Extensions []string
}

//Emitter emitter
type Emitter struct {
	Channel chan string
	Wait    time.Duration
}

// Watcher Model
type Watcher struct {
	ch            chan bool
	Options       *WatcherOption
	Fs            *fsnotify.Watcher
	WatchingSlice []string
	Emitter       map[fsnotify.Op][]*Emitter
}
