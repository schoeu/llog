package agent

import (
	"github.com/schoeu/llog/config"
	"io"
	"os"
	"time"

	"github.com/schoeu/llog/util"
)

type logInfo struct {
	Sc      *config.SingleConfig
	Status  [2]int64
	FileIns *os.File
}

func delInfo(k string) {
	li, err := getLogInfoIns(k)
	util.ErrHandler(err)
	if li != nil && li.FileIns != nil {
		err := li.FileIns.Close()
		util.ErrHandler(err)
	}
	if sm.Has(k) {
		sm.Remove(k)
	}
}

func initState(paths []string, sc config.SingleConfig) {
	seekType := getSeekType(sc)
	for _, v := range paths {
		if v != "" {
			f, offset := getFileIns(v, seekType)
			sm.SetIfAbsent(v, logInfo{
				Sc:      &sc,
				FileIns: f,
				Status:  [2]int64{offset, time.Now().Unix()},
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

func getSeekType(sc config.SingleConfig) int {
	seekType := io.SeekStart
	if sc.TailFiles {
		seekType = io.SeekEnd
	}
	return seekType
}
