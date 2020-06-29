package agent

import (
	"bufio"
	"fmt"
	"os"

	"github.com/schoeu/llog/config"
)

func logInput(sc config.SingleConfig) {
	// collect log
	fileGlob(sc)

	// close file handle schedule
	closeFileHandle(sc)

	// watch new log file schedule
	reScanTask(sc.ScanFrequency)
}

func stdInput() {
	bio := bufio.NewReader(os.Stdin)
	fmt.Println("before")
	for {
		content, _ := bio.ReadString('\n')
		fmt.Println(content)
		if content != "" {
			fmt.Println("stdin", string(content))
		}
	}
}
