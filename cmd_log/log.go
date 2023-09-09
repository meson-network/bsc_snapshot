package cmd_log

import (
	"github.com/coreservice-io/log"
	"github.com/meson-network/bsc-data-file-utils/basic"
)

func StartLog(onlyerr bool, num int64) {
	if num == 0 {
		num = 20
	}
	if onlyerr {
		basic.Logger.PrintLastN(num, []log.LogLevel{log.PanicLevel, log.FatalLevel, log.ErrorLevel})
	} else {
		basic.Logger.PrintLastN(num, []log.LogLevel{log.PanicLevel, log.FatalLevel, log.ErrorLevel, log.InfoLevel, log.WarnLevel, log.DebugLevel, log.TraceLevel})
	}
}
