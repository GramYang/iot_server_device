package log

import (
	"iot_server6/config"
	"os"

	g "github.com/GramYang/gylog"
)

var logFile = "iot_server6_log.txt"

func InitLog() {
	g.SetFlags(g.Lshortfile)
	g.SetLevel(g.LevelDebug)
	if config.Conf.LogLocal {
		g.SetOutput(os.Stderr)
	} else {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		g.SetOutput(file)
	}
	if config.Conf.LogDebug {
		g.SetLevel(g.LevelDebug)
	} else {
		g.SetLevel(g.LevelInfo)
	}
}
