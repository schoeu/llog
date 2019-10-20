package agent

import (
	"io"
	"os"
	"time"

	"github.com/schoeu/llog/util"
)

// Element: [offset, lastRead]
type logStatus map[string][2]int64

var lsCh = make(chan logStatus)
var fileCh = make(chan map[string]*os.File)
var lsCtt = logStatus{}
var delCh = make(chan string)
var timeoutDel = make(chan int)
var fileIns map[string]*os.File

func updateState() {
	defer util.Recover()

	for {
		select {
		case file := <-fileCh:
			for k, v := range file {
				if fileIns == nil {
					fileIns = map[string]*os.File{}
				}
				if fileIns[k] == nil {
					fileIns[k] = v
				}
			}
		case fileState := <-lsCh:
			for k, v := range fileState {
				lsCtt[k] = v
			}
		case k := <-delCh:
			delInfo(k)
		case aliveTime := <-timeoutDel:
			for k, v := range fileIns {
				if v != nil {
					if time.Since(time.Unix(lsCtt[k][1], 0)) > time.Second*time.Duration(aliveTime) {
						delInfo(k)
					}
				}
			}
		}
	}
}

func delInfo(k string) {
	delete(lsCtt, k)
	err := fileIns[k].Close()
	util.ErrHandler(err)
	delete(fileIns, k)
}

func initState(paths []string, sc *util.SingleConfig) {
	seekType := getSeekType(sc)
	for _, v := range paths {
		f, offset := getFileIns(v, seekType)
		fileCh <- map[string]*os.File{
			v: f,
		}
		lsCh <- logStatus{
			v: {offset, time.Now().Unix()},
		}
	}
}

func getFileIns(p string, seek int) (*os.File, int64) {
	if p != "" {
		f, err := os.Open(p)
		util.ErrHandler(err)
		offset, err := f.Seek(0, seek)
		util.ErrHandler(err)
		return f, offset
	}
	return nil, 0
}

func getSeekType(sc *util.SingleConfig) int {
	seekType := io.SeekStart
	if sc.TailFiles {
		seekType = io.SeekEnd
	}
	return seekType
}
