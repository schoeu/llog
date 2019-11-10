package agent

import (
	"bytes"
	"context"

	"github.com/schoeu/llog/util"
	"net/http"
)

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

func esPush(data *logStruct) {
	defer util.Recover()

	_, err := indexServer.BodyJson(data).
		Do(context.Background())
	// if err is not nil: es connect closed
	util.ErrHandler(err)
}
