package agent

import (
	"fmt"
	"github.com/schoeu/gopsinfo"
)

func PushData(data *gopsinfo.PsInfo, appId, secret string) {
	// TODO: 数据传输
	fmt.Println("data-> ", data.Load, appId, secret)
}
