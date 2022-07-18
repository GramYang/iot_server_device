package main

import (
	aio "iot_server6/aliyun_iot_client"
	"iot_server6/config"
	"iot_server6/log"
)

func main() {
	config.SetUp()
	log.InitLog()
	aio.BeginSubscribe()
}
