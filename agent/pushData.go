package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/schoeu/pslog_agent/util"
)

func PushData(data map[string]interface{}) {
	// TODO: 数据传输

	d, err := json.Marshal(data)
	fmt.Println(string(d))
	client := &http.Client{}
	req, err := http.NewRequest("POST", util.LogServer, bytes.NewBuffer(d))
	util.ErrHandler(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	util.ErrHandler(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.ErrHandler(err)

	fmt.Println(string(body))
}
