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
	//client := &http.Client{}
	//req, err := http.NewRequest("POST", util.LogServer, string(d))
	//util.ErrHandler(err)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//resp, err := client.Do(req)
	//util.ErrHandler(err)

	//d := `{"errorType":"resourceError","target":{"outerHTML":"<img src: \"zzzz.jpg\">","src":"http://172.23.27.232:9999/zzzz.jpg","tagName":"IMG","className":"","XPath":"/html/body/ul/li[7]/img","timeStamp":7102.404999999635},"errorLink":"http://172.23.27.232:9999/zzzz.jpg","resourceHost":"172.23.27.232","resourceType":"IMG","userBehavior":[],"id":"qOAap","sysPlat":"Other","agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36","sysVersion":"1.0","runPlat":"OTHER","phoneBrand":"测试","url":"http://172.23.27.232:9999/","referrer":"http://172.23.27.232:9999/","urlParams":"{}","network":"4g","winW":2560,"winH":1440,"locale":"zh-CN","metaData":"Document","createTime":1566637174180,"authkey":"18044033-3434-4f4b-955b-774bbf3dea64","sysMv":"1","uvId":"bBSfBekaZ3DhQ5NM","sdkVersion":"_sdkVersion_","jre":"other","cname":"","br":"chrome","brv":"73.0.3683.86","type":"error","clientTime":1566637179184}`
	//d := `{"authKey":"18044033-3434-4f4b-955b-774bbf3dea64","cpuModel":"Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz","cpuPercent":0,"dateTime":"2019-08-26T10:48:01","diskTotal":3000882061312,"diskUsed":154000150528,"diskUsedPercent":5.131829488182897,"external":"1647182","heapTotal":"6107136","heapUsed":"2890720","largeObject":"0","largeObjectUsed":"0","load":"1.96,2.12,2.09","logId":"db0cb925-577f-4a2c-994c-f4690bf544c2","logicalCores":16,"mapSpace":"528384","mapSpaceUsed":"221920","memTotal":34359738368,"memUsed":15286611968,"memUsedPercent":44.489896297454834,"newSpace":"1048576","newSpaceUsed":"565864","oldSpace":"3268608","oldSpaceUsed":"1915760","os":"darwin","physicalCores":8,"platform":"darwin","platformFamily":"Standalone Workstation","platformVersion":"10.14.5","recvSpeed":0,"rss":"26161152","sentSpeed":0,"type":"nmAgent","version":"1.0.0"}`
	resp, err := http.Post(util.LogServer, "application/json", bytes.NewBuffer(d))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.ErrHandler(err)

	fmt.Println(string(body))
}
