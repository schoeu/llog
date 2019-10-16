package agent

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/schoeu/gopsinfo"
	"github.com/schoeu/llog/util"
)

type logStruct map[string]string

func fileGlob(allLogs []string) {
	go updateState()

	for _, v := range allLogs {
		v = pathPreProcess(v)
		paths, err := filepath.Glob(v)
		util.ErrHandler(err)
		// update file state.
		initState(paths)
		watch(paths)
		tailFile(paths)
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

func tailFile(p []string) {
	seek := getSeekType()
	for _, v := range p {
		tail(v, seek)
	}
}

func tell(f *os.File, r *bufio.Reader) (offset int64, err error) {
	if f == nil {
		return
	}
	offset, err = f.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}

	if r == nil {
		return
	}

	offset -= int64(r.Buffered())
	return
}

func tail(p string, seek int) {
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

	f, _ := getFileIns(p, seek)
	r := bufio.NewReader(f)

	go func() {
		var offset int64
		var err error
		for {
			offset, err = tell(f, r)
			if err != nil {
				return
			}

			line, err := readLine(r)
			if err == nil {
				if len(include) > 0 && !util.IsInclude(line, include) {
					continue
				}
				if len(exclude) > 0 && util.IsInclude(line, exclude) {
					continue
				}

				if confMaxByte != 0 && len(line) > confMaxByte {
					line = line[:confMaxByte]
				}

				// 多行模式
				if multiline != "" {
					// 匹配开始头
					if util.IsInclude(line, []string{multiline}) {
						if logContent.Len() > 0 {
							doPush(sysInfo, st, logContent.Bytes(), apiServer)
							logContent = bytes.Buffer{}
						}
					}
					// 匹配多行其他内容
					if maxLines != 0 && logContent.Len() < maxLines {
						logContent.WriteString(line)
						continue
					}
				} else {
					doPush(sysInfo, st, []byte(line), apiServer)
				}
			} else if err == io.EOF {
				if line != "" {
					// this has the potential to never return the last line if
					// it's not followed by a newline; seems a fair trade here
					_, err := f.Seek(offset, 0)
					if err != nil {
						return
					}
				}

				select {
				case _ = <-changCh:
				}
			} else {
				// non-EOF error
				fmt.Println("Error reading", p, err)
				return
			}
		}
	}()
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return line, err
	}

	line = strings.TrimRight(line, "\n")
	return line, err
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
