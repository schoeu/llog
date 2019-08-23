package agent

import (
	"fmt"
	"github.com/schoeu/gopsinfo"
	//"io/ioutil"
	//"net/http"
	//"strings"
	//
	//"github.com/schoeu/pslog_agent/util"
)

type logData map[string]string

func PushData(data *gopsinfo.PsInfo, appId, secret,nodeInfo string) {
	// TODO: 数据传输
	fmt.Println("data-> ", data, appId, secret, nodeInfo)

	//dataProc(data, appId, secret, nodeInfo)
	//
	//client := &http.Client{}
	//
	//req, err := http.NewRequest("POST", "http:////nlogtj.zuoyebang.cc/log/performance", strings.NewReader(""))
	//util.ErrHandler(err)
	//req.Header.Set("Content-Type", "application/json")
	//
	//
	//resp, err := client.Do(req)
	//util.ErrHandler(err)
	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//util.ErrHandler(err)
	//
	//fmt.Println(string(body))
}

func dataProc(data *gopsinfo.PsInfo, appId, secret,nodeInfo string) {

}
