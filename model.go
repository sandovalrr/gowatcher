package gowatcher

import "time"

// Watcher Configuration for an instance
type Config struct {
	Dirs       []string
	Recursive  bool
	Extensions []string
	Wait       time.Duration
}
