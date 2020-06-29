package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/config"
	"github.com/schoeu/llog/util"
)

const aliveTimeDefault = 300
const freqDefault = 600

func closeFileHandle(sc config.SingleConfig) {
	aliveTime := sc.CloseInactive
	if aliveTime < 1 {
		aliveTime = aliveTimeDefault
	}
	ticker := time.NewTicker(time.Duration(aliveTime*60) * time.Second)

	go func() {
		defer util.Recover()

		for {
			<-ticker.C
			for _, v := range sm.Keys() {
				li, err := getLogInfoIns(v)
				util.ErrHandler(err)
				if li != nil && stringEqual(li.Sc.LogDir, sc.LogDir) && time.Since(time.Unix(li.Status[1], 0)) > time.Second*time.Duration(aliveTime) {
					delInfo(v)
				}
			}
		}
	}()
}

func reScanTask(freq int) {
	if freq < 1 {
		freq = freqDefault
	}
	go func() {
		defer util.Recover()

		ticker := time.NewTicker(time.Duration(freq) * time.Second)
		for {
			<-ticker.C
			reScan()
		}
	}()
}

func sysInfo(during int) {
	var psInfo gopsinfo.PsInfo
	var d time.Duration
	if during < 1 {
		d = 1
	} else if during == 0 {
		d = 10
	}
	ticker := time.NewTicker(d * time.Second)
	go func() {
		defer util.Recover()
		for {
			<-ticker.C

			psInfo = gopsinfo.GetPsInfo(d)
			sysData, err := json.Marshal(psInfo)
			util.ErrHandler(err)
			doPush(&sysData, systemType, "")
		}
	}()
}

func takeSnap(snd int) {
	ticker := time.NewTicker(time.Duration(snd) * time.Second)
	go func() {
		defer util.Recover()
		for {
			<-ticker.C

			store := storeState{}
			snap := getSnapPath()
			exist, err := util.PathExist(snap)
			util.ErrHandler(err)
			if !exist {
				err = os.Mkdir(filepath.Dir(snap), os.ModePerm)
			}
			f, err := os.Create(snap)
			for _, v := range sm.Keys() {
				li, err := getLogInfoIns(v)
				util.ErrHandler(err)
				if li != nil {
					store[v] = li.Status
				}
			}
			d, err := json.Marshal(store)
			util.ErrHandler(err)
			_, err = f.Write(d)
			err = f.Close()
			util.ErrHandler(err)
		}
	}()
}

func debugInfo() {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		defer util.Recover()
		for {
			<-ticker.C
			//for k, v := range sm.Items() {
			//	val := v.(logInfo)
			//	fmt.Println("[LLOG]", k, val.Sc)
			//}
			//fmt.Println("[LLOG]", sm.Keys())
			fmt.Fprint(os.Stdin, "xx\n")
		}
	}()
}

func getSnapPath() string {
	conf := config.GetConfig()
	snap := conf.SnapShot.SnapshotDir
	if snap == "" {
		snap = filepath.Join(util.GetTempDir(), util.SnapshotDir, util.SnapshotFile)
	}
	fmt.Println("[LLOG] snapshot file path: ", snap)
	return snap
}

func stringEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func reScan() {
	inputs := config.GetConfig().Input
	for _, v := range inputs {
		// collect log.
		fileGlob(v)
	}
}
