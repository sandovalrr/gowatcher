package gowatcher

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

// File represents a YAML configuration file that namespaces all Galley
// configuration under the top-level "clair" key.
type File struct {
	Galley Config `yaml:"galley"`
}

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

// Config is the global configuration for an instance of Galley.
type Config struct {
	Watcher  *WatcherConfig
	Database *DatabaseConfig
}

type DatabaseConfig struct {
	Port     uint16
	Host     string
	User     string
	Password string
	Database string
}

// Watcher Model
type Watcher struct {
	ch            chan bool
	Options       *WatcherOption
	Fs            *fsnotify.Watcher
	WatchingSlice []string
	Emitter       map[fsnotify.Op][]*Emitter
}
