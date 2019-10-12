package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hpcloud/tail"
	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/util"
)

type logStruct map[string]interface{}

type allLogState struct {
	offset   int64
	lastRead time.Time
	alive    bool
	tail     *tail.Tail
}

var (
	allPath = map[string]allLogState{}
)

func fileGlob() {
	excludeFiles := util.GetConfig().ExcludeFiles
	for _, v := range allLogs {
		v = pathPreProcess(v)
		paths, err := filepath.Glob(v)

		util.ErrHandler(err)
		for _, v := range paths {
			// log path store.
			if !allPath[v].alive {
				if len(excludeFiles) > 0 && util.IsInclude(v, excludeFiles) {
					continue
				}
				fmt.Println("watch new file: ", v)
				go logFilter(v)
			}
		}
	}
}

func pathPreProcess(p string) string {
	exist, err := util.PathExist(p)
	util.ErrHandler(err)
	if !exist {
		err = os.Mkdir(filepath.Dir(p), os.ModePerm)
	}
	if !filepath.IsAbs(p) {
		p = util.GetAbsPath("", p)
	}
	return p
}

func logFilter(logFile string) {
	conf := util.GetConfig()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	seekType := io.SeekStart
	if conf.TailFiles {
		seekType = io.SeekEnd
	}
	t, err := tail.TailFile(logFile, tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: seekType,
		},
		Follow: true,
	})

	//allPath[logFile] = allLogState{
	//	tail: t,
	//	alive: true,
	//}

	util.ErrHandler(err)

	st := time.Now()
	var logContent bytes.Buffer

	include, exclude, apiServer, multiline := conf.Include, conf.Exclude, conf.ApiServer, conf.Multiline.Pattern
	sysInfo, confMaxByte, maxLines := conf.SysInfo, conf.MaxBytes, conf.Multiline.MaxLines
	for line := range t.Lines {
		//offset, _ := t.Tell()
		//allPath[logFile] = offset
		//allPath[logFile] = allLogState{
		//	lastRead: time.Now(),
		//}

		text := line.Text
		if len(include) > 0 && !util.IsInclude(text, include) {
			continue
		}
		if len(exclude) > 0 && util.IsInclude(text, exclude) {
			continue
		}

		if confMaxByte != 0 && len(text) > confMaxByte {
			text = text[:confMaxByte]
		}

		// 多行模式
		if multiline != "" {
			// 匹配开始头
			if util.IsInclude(text, []string{multiline}) {
				if logContent.Len() > 0 {
					doPush(sysInfo, st, logContent.Bytes(), apiServer)
					logContent = bytes.Buffer{}
				}
			}
			// 匹配多行其他内容
			if maxLines != 0 && logContent.Len() < maxLines {
				logContent.WriteString(text)
				continue
			}
		} else {
			doPush(sysInfo, st, []byte(text), apiServer)
		}
	}
	util.ErrHandler(err)
}

func doPush(sysInfo bool, st time.Time, text []byte, apiServer string) {
	var rs logStruct

	if sysInfo {
		var psInfo gopsinfo.PsInfo
		et := time.Now()
		during := et.Sub(st)
		timeSub := int(during)
		if timeSub < 1 {
			during = time.Microsecond * 1000
		}
		psInfo = gopsinfo.GetPsInfo(during)
		st = et

		sysData, err := json.Marshal(psInfo)
		util.ErrHandler(err)
		rs = util.CombineData(logStruct{
			"@sysInfo": string(sysData),
			"@message": string(text),
		})
	} else {
		rs = logStruct{
			"@message": string(text),
		}
	}

	if apiServer != "" {
		go pushData(combineTags(rs), apiServer)
	}
}

func combineTags(rs logStruct) logStruct {
	// 日志签名
	rs["@version"] = util.Version
	rs["@logId"] = util.UUID()
	rs["@type"] = util.AppName
	rs["@timestamps"] = time.Now().UnixNano() / 1e6
	return rs
}

func closeFileHandle() {
	fmt.Println("closeFileHandle")
	conf := util.GetConfig()
	aliveTime := conf.CloseInactive
	for key, v := range allPath {
		if !v.alive {
			if time.Since(allPath[key].lastRead) > time.Second*time.Duration(aliveTime) {
				tailErr := v.tail.Stop()
				delete(allPath, key)
				util.ErrHandler(tailErr)
				// delete map data.
			}
		}
	}
}
