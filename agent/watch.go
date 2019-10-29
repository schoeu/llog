package agent

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/schoeu/llog/util"
)

var once sync.Once
var fsWatcher *fsnotify.Watcher

func addWatch(paths []string, sc *util.SingleConfig) {
	var err error
	once.Do(func() {
		fsWatcher, err = fsnotify.NewWatcher()
		util.ErrHandler(err)
	})

	//TODO: defer fsWatcher.Close()
	excludeFiles := sc.ExcludeFiles
	for _, v := range paths {
		// log path store.
		if len(excludeFiles) > 0 && util.IsInclude([]byte(v), excludeFiles) {
			continue
		}
		fmt.Println("watch file: ", v)
		err = fsWatcher.Add(v)
		util.ErrHandler(err)
	}
}

func watch(sc *util.SingleConfig) {
	defer util.Recover()

	var push = lineFilter(sc)
	for {
		select {
		case ev := <-fsWatcher.Events:
			// add new file
			if ev.Op&fsnotify.Create == fsnotify.Create {
				if ev.Name != "" {
					initState([]string{ev.Name}, sc)
					err := fsWatcher.Add(ev.Name)
					util.ErrHandler(err)
				}
			}
			//change file content
			if ev.Op&fsnotify.Write == fsnotify.Write {
				key := ev.Name
				fmt.Println("change ->", key)
				if key != "" {
					f := fileIns[key]
					var count int
					offset, err := f.Seek(0, io.SeekCurrent)
					util.ErrHandler(err)
					line := bufio.NewReader(f)
					var content []byte
					for {
						content, _, err = line.ReadLine()
						if err == io.EOF {
							break
						}
						count += len(content)
						push(&content)
					}
					if err == io.EOF {
						//_, seekErr := f.Seek(offset, io.SeekStart)
						lsCh <- logStatus{
							key: {offset + int64(count+1), time.Now().Unix()},
						}
						continue
					}
				}
			}
			// remove log file
			if ev.Op&fsnotify.Remove == fsnotify.Remove {
				if ev.Name != "" {
					delCh <- ev.Name
				}
			}
			// rename log file
			if ev.Op&fsnotify.Rename == fsnotify.Rename {
				if ev.Name != "" {
					delCh <- ev.Name
					initState([]string{ev.Name}, sc)
				}
			}
		case err := <-fsWatcher.Errors:
			if err != io.EOF {
				util.ErrHandler(err)
			}
		}
	}
}
