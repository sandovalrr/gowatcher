package gowatcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/pkg/capnslog"
	"github.com/fsnotify/fsnotify"
	"github.com/sandovalrr/gowatcher/utils"
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

	go watcher.handleEvents()

	for _, file := range watcher.Options.Dirs {
		watcher.watch(file)
	}

	<-done
}

func (watcher *Watcher) watch(path string) {
	isDirectory := utils.IsDirectory(path)

	if !isDirectory {
		log.Infof("%s is not directory", path)
		return
	}

	if utils.SliceStringContains(watcher.WatchingSlice, path) {
		log.Infof("%s is being watched", path)
		return
	}

	err := os.MkdirAll(path, 0777)
	if err != nil {
		log.Error(err)
		return
	}

	err = watcher.Fs.Add(path)
	if err != nil {
		log.Error(err)
		return
	}

	watcher.WatchingSlice = append(watcher.WatchingSlice, path)
	log.Infof("%s append to Watcher\n", path)

	if watcher.Options.Recursive {
		for _, subDir := range utils.GetSubFolders(path) {
			watcher.watch(subDir)
		}
	}

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
	isDirectory := utils.IsDirectory(event.Name)

	if !isDirectory {
		if ext, ok := watcher.isExtensionFileValid(event); !ok {
			log.Infof("Ignoring file: %s, extension %s is not able to be processed.", event.Name, ext)
			return
		}
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		if isDirectory {
			if utils.SliceStringContains(watcher.WatchingSlice, event.Name) {
				watcher.Fs.Remove(event.Name)
				watcher.WatchingSlice = utils.SliceRemoveString(watcher.WatchingSlice, event.Name)
			}
		}
	}

	if event.Op&fsnotify.Create == fsnotify.Create {

		if isDirectory {
			log.Infof("New Directory %s created, checking if recursive flag is on to start watching files", event.Name)
			return
		}
	}

	watcher.emit(event)
}

func (watcher *Watcher) emit(event *fsnotify.Event) {
	if emitters, ok := watcher.Emitter[event.Op]; ok {
		for _, emitter := range emitters {

			go func(emitter *Emitter, event *fsnotify.Event) {
				select {
				case <-time.After(emitter.Wait * time.Second):
					log.Infof("Trigering %s event on file %s", event.Op.String(), event.Name)
					emitter.Channel <- event.Name
				}
			}(emitter, event)

		}
	}
}

//Subscribe subscribe to events
func (watcher *Watcher) Subscribe(op fsnotify.Op, emitter *Emitter) {

	if watcher.Emitter == nil {
		watcher.Emitter = make(map[fsnotify.Op][]*Emitter)
	}

	if _, ok := watcher.Emitter[op]; ok {
		watcher.Emitter[op] = append(watcher.Emitter[op], emitter)
	} else {
		watcher.Emitter[op] = []*Emitter{emitter}
	}
}

//UnSubscribe unsubscribe to events
func (watcher *Watcher) UnSubscribe(op fsnotify.Op, emitter *Emitter) {

	if _, ok := watcher.Emitter[op]; ok {
		for index, e := range watcher.Emitter[op] {
			if e == emitter {
				watcher.Emitter[op] = append(watcher.Emitter[op][:index], watcher.Emitter[op][index+1:]...)
			}
		}
	}
}
