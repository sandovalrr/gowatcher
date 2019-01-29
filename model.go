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
	Wait       time.Duration
}

// Watcher Model
type Watcher struct {
	Options *WatcherOption
	Fs      *fsnotify.Watcher
}
