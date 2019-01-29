package gowatcher

import (
	"path/filepath"
	"strings"

	"github.com/coreos/pkg/capnslog"
	"github.com/fsnotify/fsnotify"
)

var log = capnslog.NewPackageLogger(Repo, "watcher")

// NewWatcher new watcher
func NewWatcher(options WatcherOption) *Watcher {
	return &Watcher{
		Options: &options,
	}
}

//Start start
func (watcher *Watcher) Start() {

	if len(watcher.Options.Dirs) == 0 {
		log.Info("Watcher Disabled")
		return
	}

	fs, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return
	}

	watcher.Fs = fs

	defer fs.Close()

	done := make(chan bool)

	<-done
}

func (watcher *Watcher) handleEvents() {
	for {
		select {
		case event := <-watcher.Fs.Events:
			watcher.onEvent(&event)
		case err := <-watcher.Fs.Errors:
			log.Error(err)
		}
	}
}

func (watcher *Watcher) isExtensionFileValid(event *fsnotify.Event) (string, bool) {

	if len(watcher.Options.Extensions) == 0 {
		return "*", true
	}

	ext := filepath.Ext(event.Name)
	isValid := false

	for _, extension := range watcher.Options.Extensions {
		isValid = isValid || strings.ToLower(ext) == strings.ToLower(extension)

		if isValid {
			break
		}
	}

	return ext, isValid
}

func (watcher *Watcher) onEvent(event *fsnotify.Event) {

	if ext, ok := watcher.isExtensionFileValid(event); !ok {
		log.Infof("Ignoring file: %s, extension %s is not able to be processed.", event.Name, ext)
		return
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		watcher.onRemove(event)
		return
	}

	if event.Op&fsnotify.Create == fsnotify.Create {
		watcher.onCreate(event)
		return
	}

}

func (watcher *Watcher) onCreate(event *fsnotify.Event) {
	log.Infof("%s created. %v seconds to trigger event", event.Name, watcher.Options.Wait.Seconds)
}

func (watcher *Watcher) onRemove(event *fsnotify.Event) {
	log.Infof("%s removed", event.Name)
}
