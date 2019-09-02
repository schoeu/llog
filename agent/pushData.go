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
	resp, err := http.Post(util.LogServer, "application/json", bytes.NewBuffer(d))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.ErrHandler(err)

	fmt.Println("log server:\n", string(body))
}
