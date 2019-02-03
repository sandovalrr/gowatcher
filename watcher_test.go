package gowatcher_test

import (
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/sandovalrr/gowatcher"
)

var dirs = []string{"./assets"}
var exts = []string{".mp3"}

func TestNewWatcher(t *testing.T) {
	watcher := gowatcher.NewWatcher(gowatcher.WatcherOption{
		Dirs:       dirs,
		Extensions: exts,
		Recursive:  false,
	})

	if watcher == nil {
		t.Error("Expected watcher were defined")
	}
}

func TestStart(t *testing.T) {
	watcher := gowatcher.NewWatcher(gowatcher.WatcherOption{
		Dirs:       dirs,
		Extensions: exts,
		Recursive:  false,
	})

	go watcher.Start()
	time.Sleep(1 * time.Second)

	if !watcher.IsWatching() {
		t.Error("Watcher is not running")
	}
	watcher.Close()
}

func TestStop(t *testing.T) {
	watcher := gowatcher.NewWatcher(gowatcher.WatcherOption{
		Dirs:       dirs,
		Extensions: exts,
		Recursive:  false,
	})

	go watcher.Start()
	time.Sleep(1 * time.Second)
	watcher.Close()

	if watcher.IsWatching() {
		t.Error("Watcher still running")
	}
}

func TestSubscribe(t *testing.T) {
	watcher := gowatcher.NewWatcher(gowatcher.WatcherOption{
		Dirs:       dirs,
		Extensions: exts,
		Recursive:  false,
	})

	go watcher.Start()
	subscriber := &gowatcher.Emitter{
		Channel: make(chan string),
		Wait:    time.Duration(0),
	}
	watcher.Subscribe(fsnotify.Create, subscriber)

	if len(watcher.Emitter[fsnotify.Create]) == 0 {
		t.Error("No subscribed to Create event")
	}
}

func TestUnSubscribe(t *testing.T) {
	watcher := gowatcher.NewWatcher(gowatcher.WatcherOption{
		Dirs:       dirs,
		Extensions: exts,
		Recursive:  false,
	})

	go watcher.Start()
	subscriber := &gowatcher.Emitter{
		Channel: make(chan string),
		Wait:    time.Duration(0),
	}
	watcher.Subscribe(fsnotify.Create, subscriber)
	watcher.UnSubscribe(fsnotify.Create, subscriber)

	emitters := watcher.Emitter[fsnotify.Create]
	if len(emitters) > 0 {
		t.Error("Subscriber still subscribed to Create event")
	}
}
