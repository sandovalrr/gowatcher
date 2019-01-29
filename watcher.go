package gowatcher

// NewWatcher new watcher
func NewWatcher(options WatcherOption) *Watcher {
	return &Watcher{
		Options: &options,
	}
}

func (watcher *Watcher) start() {

	if len(watcher.Options.Dirs) == 0 {
		//TODO: watcher service disabled
		return
	}

}
