package agent

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/schoeu/llog/util"
)

type logStruct map[string]string

var apiServer, name string

const errorType = "error"
const normalType = "normal"
const systemType = "system"

func fileGlob(sc *util.SingleConfig) {
	allLogs := sc.LogDir
	if len(allLogs) == 0 {
		logFileDir := util.GetTempDir()
		allLogs = append(allLogs, filepath.Join(logFileDir, util.LogDir, util.FilePattern))
	}
	// allLogs: - /var/logs/**/*.log
	for _, v := range allLogs {
		v = pathPreProcess(v)
		// paths: ["/var/logs/1.log","/var/logs/2.log"]
		paths, err := filepath.Glob(v)
		util.ErrHandler(err)
		// update file state.
		initState(paths, sc)
		addWatch(paths, sc)
	}

	go watch(sc)
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

	if maxLines == 0 {
		maxLines = 10
	}

	var lineCount int
	return func(l *[]byte) {
		line := *l

		// 多行模式
		if multiline != "" {
			// 匹配开始头
			if util.IsInclude(line, []string{multiline}) {
				if logContent.Len() > 0 {
					ok, rs := filter(include, exclude, logContent.Bytes(), confMaxByte)
					if ok {
						return
					}
					doPush(rs, errorType)
					logContent = bytes.Buffer{}
					lineCount = 0
				}
			}
			lineCount++
			// 匹配多行其他内容
			if lineCount < maxLines {
				logContent.Write(line)
			}
		} else {
			ok, rs := filter(include, exclude, line, confMaxByte)
			if ok {
				return
			}
			doPush(rs, normalType)
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
	return false, &line
}

func doPush(text *[]byte, types string) {
	// 日志签名
	var rs = logStruct{
		"@message":    string(*text),
		"@version":    util.Version,
		"@logId":      util.UUID(),
		"@timestamps": strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		"@types":      types,
		"@name":       name,
	}

	if apiServer != "" {
		go apiPush(&rs, apiServer)
	}

	if indexServer != nil {
		go esPush(&rs)
	}
}
