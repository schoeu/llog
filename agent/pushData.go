package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/schoeu/nma/util"
)

func PushData(data map[string]interface{}) {
	d, err := json.Marshal(data)
	fmt.Println(string(d))
	client := &http.Client{}
	req, err := http.NewRequest("POST", util.LogServer, bytes.NewBuffer(d))
	util.ErrHandler(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://h5.zuoyebang.cc")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")

	resp, err := client.Do(req)
	util.ErrHandler(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.ErrHandler(err)

	fmt.Println(string(body))
}
