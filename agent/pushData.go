package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/schoeu/nma/util"
)

func PushData(data map[string]interface{}, server string) {
	d, err := json.Marshal(data)

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, time.Second*2)
				if err != nil {
					return nil, err
				}
				err = conn.SetDeadline(time.Now().Add(time.Second * 90))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
	request, err := http.NewRequest("POST", server, bytes.NewBuffer(d))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
		return
	}
	resp, err := client.Do(request)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.ErrHandler(err)
	fmt.Println("log server response:\n", string(body))
}
