package agent

import (
	"io"
	"os"
	"time"

	"github.com/schoeu/llog/util"
)

func delInfo(k string) {
	li := getLogInfoIns(k)
	if li.fileIns != nil {
		err := li.fileIns.Close()
		util.ErrHandler(err)
	}
	sm.Remove(k)
}

func initState(paths []string, sc *util.SingleConfig) {
	seekType := getSeekType(sc)
	for _, v := range paths {
		if v != "" {
			f, offset := getFileIns(v, seekType)
			sm.SetIfAbsent(v, logInfo{
				sc:      sc,
				fileIns: f,
				status:  [2]int64{offset, time.Now().Unix()},
			})
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
