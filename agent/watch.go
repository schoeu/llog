package agent

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/schoeu/llog/util"
)

var once sync.Once
var fsWatcher *fsnotify.Watcher

type logInfo struct {
	sc      *util.SingleConfig
	status  [2]int64
	fileIns *os.File
}

func addWatch() {
	var err error
	once.Do(func() {
		fsWatcher, err = fsnotify.NewWatcher()
		util.ErrHandler(err)
	})

	//TODO: defer fsWatcher.Close()
	for _, v := range sm.Keys() {
		li, err := getLogInfoIns(v)
		util.ErrHandler(err)
		if li != nil {
			excludeFiles := li.sc.ExcludeFiles
			// log path store.
			if len(excludeFiles) > 0 && util.IsInclude([]byte(v), excludeFiles) {
				continue
			}
			err := fsWatcher.Add(v)
			util.ErrHandler(err)
			fmt.Println("watch file: ", v)
		}
	}
}

func watch() {
	go func() {
		defer util.Recover()

		for {
			select {
			case ev := <-fsWatcher.Events:
				// add new file
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if ev.Name != "" {
						reScan()
					}
				}
				//change file content
				if ev.Op&fsnotify.Write == fsnotify.Write {
					key := ev.Name
					if key != "" {
						fi, err := getLogInfoIns(key)
						util.ErrHandler(err)

						if fi != nil {
							var push = lineFilter(key)
							f := fi.fileIns
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
								if push != nil {
									push(&content)
								}
							}
							if err == io.EOF {
								//_, seekErr := f.Seek(offset, io.SeekStart)
								//util.ErrHandler(seekErr)
								sm.Set(key, logInfo{
									sc:      fi.sc,
									status:  [2]int64{offset + int64(count+1), time.Now().Unix()},
									fileIns: f,
								})
								continue
							}
						}
					}
				}
				// remove log file
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if ev.Name != "" {
						delInfo(ev.Name)
					}
				}
				// rename log file
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					if ev.Name != "" {
						//delCh <- ev.Name
						//initState([]string{ev.Name}, sc)
					}
				}
			case err := <-fsWatcher.Errors:
				if err != io.EOF {
					util.ErrHandler(err)
				}
			}
		}
	}()
}
