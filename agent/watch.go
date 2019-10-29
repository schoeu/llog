package agent

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/schoeu/llog/util"
)

var once sync.Once
var sm = cmap.New()
var fsWatcher *fsnotify.Watcher

type logInfo struct {
	data    bytes.Buffer
	sc      *util.SingleConfig
	status  [2]int64
	fileIns *os.File
	lineCount int
}

func addWatch() {
	var err error
	once.Do(func() {
		fsWatcher, err = fsnotify.NewWatcher()
		util.ErrHandler(err)
	})

	//TODO: defer fsWatcher.Close()
	for _, v := range sm.Keys() {
		li := getLogInfoIns(v)
		excludeFiles := li.sc.ExcludeFiles
		// log path store.
		if len(excludeFiles) > 0 && util.IsInclude([]byte(v), excludeFiles) {
			continue
		}
		fmt.Println("watch file: ", v)
		err = fsWatcher.Add(v)
		util.ErrHandler(err)
	}
}

func watch() {
	defer util.Recover()
	go func() {
		for {
			select {
			case ev := <-fsWatcher.Events:
				// add new file
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if ev.Name != "" {
						//initState([]string{ev.Name}, sc)
						//err := fsWatcher.Add(ev.Name)
						//util.ErrHandler(err)
					}
				}
				//change file content
				if ev.Op&fsnotify.Write == fsnotify.Write {
					key := ev.Name
					if key != "" {
						var push = lineFilter(key)
						fi := getLogInfoIns(key)
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
							push(&content)
						}
						if err == io.EOF {
							//_, seekErr := f.Seek(offset, io.SeekStart)
							//lsCh <- logStatus{
							//	key: {offset + int64(count+1), time.Now().Unix()},
							//}
							fi.status = [2]int64{offset + int64(count+1), time.Now().Unix()}
							continue
						}
					}
				}
				// remove log file
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if ev.Name != "" {
						//delCh <- ev.Name
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
