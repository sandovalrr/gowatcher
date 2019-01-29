package gowatcher

import "time"

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
}
