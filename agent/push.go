package agent

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/schoeu/llog/util"
)

var (
	hasClient bool
)

func getClint() http.Client {
	client := http.Client{
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
	hasClient = true
	return client
}

func pushData(data logStruct, server string) {
	d, err := json.Marshal(data)
	var client http.Client
	if !hasClient {
		client = getClint()
	}

	request, err := http.NewRequest("POST", server, bytes.NewBuffer(d))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
		return
	}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	//_, err = ioutil.ReadAll(resp.Body)
	util.ErrHandler(err)
	//fmt.Println("log server response:\n", string(body))
}
