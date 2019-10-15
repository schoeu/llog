package agent

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/util"
)

type logStruct map[string]string

func fileGlob(allLogs []string) {
	for _, v := range allLogs {
		v = pathPreProcess(v)
		paths, err := filepath.Glob(v)
		util.ErrHandler(err)

		// update file state.
		initState(paths)
		watch(paths)
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

func tail(f *os.File, info [2]int64) {
	//f.ReadAt()
	reader := bufio.NewReader(f)
	reader.ReadLine()
	r := bufio.NewReader(f)
}

func logFilter(path []string) {
	conf := util.GetConfig()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	st := time.Now()
	var logContent bytes.Buffer

	include, exclude, apiEnable, multiline := conf.Include, conf.Exclude, conf.ApiServer.Enable, conf.Multiline.Pattern
	sysInfo, confMaxByte, maxLines := conf.SysInfo, conf.MaxBytes, conf.Multiline.MaxLines

	var apiServer string
	if apiEnable && conf.ApiServer.Url != "" {
		apiServer = conf.ApiServer.Url
	}

	for line := range t.Lines {
		offset, _ := t.Tell()
		lsCh <- logStatus{
			logFile: {offset, time.Now().Unix()},
		}

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
	var rs = logStruct{
		"@message": string(text),
	}
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
		rs["@sysInfo"] = string(sysData)
	}
	combineData := combineTags(rs)
	if apiServer != "" {
		go apiPush(combineData, apiServer)
	}

	if indexServer != nil {
		go esPush(combineData)
	}
}

func combineTags(rs logStruct) logStruct {
	// 日志签名
	rs["@version"] = util.Version
	rs["@logId"] = util.UUID()
	rs["@type"] = util.AppName
	rs["@timestamps"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	return rs
}
