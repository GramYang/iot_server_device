package main

import (
	"iot_server8/config"
	"iot_server8/log"
	"iot_server8/router"
	sc "iot_server8/sqlx_client"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	config.SetUp()
	log.InitLog()
	sc.SetUp()
	r := router.NewRouter()
	if config.Conf.LogDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	_ = r.Run(":" + strconv.Itoa(config.Conf.ServerPort))
}
