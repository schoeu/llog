package agent

import (
	"bytes"
	"context"
	"github.com/schoeu/llog/config"
	"net"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/schoeu/llog/util"
)

var (
	client      *http.Client
	indexServer *elastic.IndexService
	esIndex     string
)

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
	return client
}

func apiPush(data *logStruct) {
	defer util.Recover()

	d, err := json.Marshal(data)
	if client == nil {
		client = getClint()
	}

	request, err := http.NewRequest("POST", apiServer, bytes.NewBuffer(d))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
		return
	}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	util.ErrHandler(err)
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
	if esIndex != "" {
		indexServer = client.Index().Index(esIndex)
	}
}

func esPush(data *logStruct) {
	defer util.Recover()
	if indexServer != nil {
		_, err := indexServer.BodyJson(data).
			Do(context.Background())
		// if err is not nil: es connect closed
		util.ErrHandler(err)

	}
}
