package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/schoeu/nma/util"
)

func PushData(data map[string]interface{}, server string) {
	d, err := json.Marshal(data)
	resp, err := http.Post(server, "application/json", bytes.NewBuffer(d))

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.ErrHandler(err)
	fmt.Println("log server response:\n", string(body))
}
