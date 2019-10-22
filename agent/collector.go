package agent

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/schoeu/llog/util"
)

type logStruct map[string]string

var apiServer, name string

func fileGlob(sc *util.SingleConfig) {
	allLogs := sc.LogDir
	if len(allLogs) == 0 {
		logFileDir := util.GetTempDir()
		allLogs = append(allLogs, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
	}
	for _, v := range allLogs {
		v = pathPreProcess(v)
		paths, err := filepath.Glob(v)
		util.ErrHandler(err)
		// update file state.
		initState(paths, sc)
		watch(paths, sc)
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

func lineFilter(sc *util.SingleConfig) func(*[]byte) {
	conf := util.GetConfig()
	output := conf.Output

	var logContent bytes.Buffer

	include, exclude, apiEnable, multiline := sc.Include, sc.Exclude, output.ApiServer.Enable, sc.Multiline.Pattern
	confMaxByte, maxLines, appName := sc.MaxBytes, sc.Multiline.MaxLines, conf.Name

	if apiEnable && output.ApiServer.Url != "" {
		apiServer = output.ApiServer.Url
	}

	if appName == "" {
		appName = util.AppName
	}
	name = appName
	// flag 0: nil  1: start  2:end
	var flag bool

	return func(l *[]byte) {
		line := *l

		// 多行模式
		if multiline != "" {
			// 匹配开始头
			if util.IsInclude(line, []string{multiline}) {
				flag = !flag
				if logContent.Len() > 0 {
					ok, rs := filter(include, exclude, line, confMaxByte)
					if ok {
						return
					}
					doPush(rs, false)
					logContent = bytes.Buffer{}
				}
			}
			// 匹配多行其他内容
			if maxLines != 0 && logContent.Len() < maxLines {
				logContent.Write(line)
			}
		} else {
			ok, rs := filter(include, exclude, line, confMaxByte)
			if ok {
				return
			}
			doPush(rs, false)
		}
	}
}

func filter(include, exclude []string, line []byte, max int) (bool, *[]byte) {
	if len(include) > 0 && !util.IsInclude(line, include) {
		return true, nil
	}
	if len(exclude) > 0 && util.IsInclude(line, exclude) {
		return true, nil
	}
	if max != 0 && len(line) > max {
		line = line[:max]
	}
	fmt.Println(string(line))
	return false, &line
}

func doPush(text *[]byte, isSys bool) {
	cd := combineData(text, isSys)
	if apiServer != "" {
		go apiPush(cd, apiServer)
	}

	if indexServer != nil {
		go esPush(cd)
	}
}

func combineData(text *[]byte, isSys bool) *logStruct {
	var rs = logStruct{
		"@message": string(*text),
		// 日志签名
		"@version":    util.Version,
		"@logId":      util.UUID(),
		"@timestamps": strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		"@sysInfo":    strconv.FormatBool(isSys),
	}
	return &rs
}
