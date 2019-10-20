package agent

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/schoeu/llog/util"
)

func watch(paths []string, sc *util.SingleConfig) {
	w, err := fsnotify.NewWatcher()
	util.ErrHandler(err)
	// TODO
	//defer w.Close()
	var push = lineFilter(sc)

	excludeFiles := sc.ExcludeFiles
	for _, v := range paths {
		// log path store.
		if len(excludeFiles) > 0 && util.IsInclude([]byte(v), excludeFiles) {
			continue
		}
		fmt.Println("watch file: ", v)
		err = w.Add(v)
		util.ErrHandler(err)
	}
	go func() {
		defer util.Recover()

		for {
			select {
			case ev := <-w.Events:
				// add new file
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if ev.Name != "" {
						initState([]string{ev.Name}, sc)
						err = w.Add(ev.Name)
						util.ErrHandler(err)
					}
				}
				//change file content
				if ev.Op&fsnotify.Write == fsnotify.Write {
					key := ev.Name
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
			case err := <-w.Errors:
				if err != io.EOF {
					util.ErrHandler(err)
				}
			}
		}
	}()
}
