package log

import (
	"iot_server8/config"
	"os"

	g "github.com/GramYang/gylog"
	"github.com/gin-gonic/gin"
)

var logFile = "iot_server8_log.txt"

func InitLog() {
	g.SetFlags(g.Lshortfile | g.Lmicroseconds | g.Ldate)
	g.SetLevel(g.LevelDebug)
	if config.Conf.LogLocal {
		g.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		g.SetOutput(file)
		gin.DefaultWriter = file
	}
	if config.Conf.LogDebug {
		g.SetLevel(g.LevelDebug)
	} else {
		g.SetLevel(g.LevelInfo)
	}
}
