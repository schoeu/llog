package agent

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/schoeu/llog/config"
	"github.com/schoeu/llog/util"
)

var (
	client      *http.Client
	indexServer *elastic.IndexService
	esIndex     string
)

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

func doPush(text *[]byte, types, fields string) {
	// log signature
	var rs = logStruct{
		"@message":    string(*text),
		"@version":    util.Version,
		"@logId":      util.UUID(),
		"@timestamps": strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		"@types":      types,
		"@name":       name,
		"@fields":     "",
	}

	if fields != "" {
		rs["@fields"] = fields
	}

	if apiServer != "" {
		go apiPush(&rs)
	}
	if indexServer != nil {
		go esPush(&rs)
	}
}

func getLogInfoIns(p string) (*logInfo, error) {
	logContent, ok := sm.Get(p)
	if !ok {
		return nil, errors.New(fmt.Sprintf("[LLOG] error: %s is not exist in sync map", p))
	}
	li := logContent.(logInfo)
	return &li, nil
}

// http server api push.
func getClint() *http.Client {
	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, time.Second*2)
				if err != nil {
					return nil, err
				}
				err = conn.SetDeadline(time.Now().Add(time.Second * 90))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 90,
		},
	}
	fmt.Printf("[LLOG] API server have been initialized, API: %s\n", apiServer)
	return client
}

// es push.
func esInit() {
	output := config.GetConfig().Output
	esConf := output.Elasticsearch
	client, err := elastic.NewClient(
		elastic.SetURL(esConf.Host...),
		elastic.SetBasicAuth(esConf.Username, esConf.Password),
		elastic.SetSniff(false),
	)
	util.ErrHandler(err)
	esIndex = esConf.Index
	indexServer = client.Index()
	if esIndex != "" {
		indexServer = indexServer.Index(esIndex)
	}
	fmt.Printf("[LLOG] Elasticsearch client have been initialized, host: %s, index: %s\n", esConf.Host, esIndex)
}
