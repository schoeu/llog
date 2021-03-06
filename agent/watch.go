package agent

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	fs "github.com/fsnotify/fsnotify"
	"github.com/schoeu/llog/util"
)

func addWatchFile() *fs.Watcher {
	fsWatcher, err := fs.NewWatcher()
	util.ErrHandler(err)

	//TODO: defer fsWatcher.Close()
	if len(sm.Keys()) == 0 {
		fmt.Println("[LLOG] no file need to watch, please check config file.")
		os.Exit(1)
	} else {
		for _, v := range sm.Keys() {
			li, err := getLogInfoIns(v)
			util.ErrHandler(err)
			if li != nil {
				excludeFiles := li.Sc.ExcludeFiles
				// log path store.
				if len(excludeFiles) > 0 && util.IsInclude([]byte(v), excludeFiles) {
					continue
				}
				err := fsWatcher.Add(v)
				util.ErrHandler(err)
				fmt.Println("[LLOG] watch file: ", v)
			}
		}
	}

	return fsWatcher
}

func watch(fsWatcher *fs.Watcher) {
	defer util.Recover()

	for {
		select {
		case ev := <-fsWatcher.Events:
			if ev.Op&fs.Write == fs.Write {
				key := ev.Name
				if key != "" {
					fi, err := getLogInfoIns(key)
					util.ErrHandler(err)

					if fi != nil {
						var push = lineFilter(key)
						f := fi.FileIns
						var count int
						offset, err := f.Seek(0, io.SeekCurrent)
						util.ErrHandler(err)
						line := bufio.NewReader(f)
						for {
							content, _, err := line.ReadLine()
							if err == io.EOF {
								break
							}
							count += len(content)
							if push != nil {
								push(&content)
							}
						}
						if err == io.EOF {
							sm.Set(key, logInfo{
								Sc:      fi.Sc,
								Status:  [2]int64{offset + int64(count+1), time.Now().Unix()},
								FileIns: f,
							})
							continue
						}
					}
				}
			}
			// remove log file
			if ev.Op&fs.Remove == fs.Remove {
				if ev.Name != "" {
					delInfo(ev.Name)
				}
			}
			// rename log file
			if ev.Op&fs.Rename == fs.Rename {
				if ev.Name != "" {
					reScan()
				}
			}
		case err := <-fsWatcher.Errors:
			if err != io.EOF {
				util.ErrHandler(err)
			}
		}
	}
}
