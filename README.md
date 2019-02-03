<h1 align="center">Go Watcher</h1>
Light go library to watch directory files

## Instalation

```text
$ go get github.com/sandovalrr/gowatcher
```

or using glide

```text
$ glide get github.com/sandovalrr/gowatcher
```

## Usage

```go

import (

  "time"

  "github.com/fsnotify/fsnotify"

  "github.com/sandovalrr/gowatcher"
)

//...
//...

watcher := gowatcher.NewWatcher(gowatcher.WatcherOption{
  Dirs: []string{"test_path_to_watch","another_path_to_watch"},
  Recursive: true,
  Extensions: []string{".mp3",".wav"},
})

go watcher.Start()

subscriber := &gowatcher.Emitter{
  Channel: make(chan string),
  Wait:    time.Duration(0),
}

watcher.Subscribe(fsnotify.Create, subscriber)

//...
go func(subscriber *gowatcher.Emitter){
  for {
    select {
      case path := <- subscriber.Channel:
      //Action on event
    }
  }
}(subscriber)
//...

//onFinish
watcher.Close()

```

## API
