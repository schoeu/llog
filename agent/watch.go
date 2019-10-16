package agent

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/schoeu/llog/util"
)

func watch(paths []string) {
	w, err := fsnotify.NewWatcher()
	util.ErrHandler(err)
	// TODO
	//defer w.Close()

	excludeFiles := util.GetConfig().ExcludeFiles
	for _, v := range paths {
		// log path store.
		if len(excludeFiles) > 0 && util.IsInclude(v, excludeFiles) {
			continue
		}
		fmt.Println("watch file: ", v)
		err = w.Add(v)
		util.ErrHandler(err)
	}

	go func() {
		for {
			select {
			case ev := <-w.Events:
				// add new file
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if ev.Name != "" {
						initState([]string{ev.Name})
					}
				}
				//change file content
				if ev.Op&fsnotify.Write == fsnotify.Write {
					if ev.Name != "" {
						changCh <- ev.Name
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
						initState([]string{ev.Name})
					}
				}
			case err := <-w.Errors:
				util.ErrHandler(err)
				return
			}
		}
	}()
}
