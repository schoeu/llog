package agent

import "github.com/schoeu/llog/config"

func logInput(sc config.SingleConfig) {
	// collect log
	fileGlob(sc)

	// close file handle schedule
	closeFileHandle(sc)

	// watch new log file schedule
	reScanTask(sc.ScanFrequency)
}

func stdInput(sc config.SingleConfig) {

}
