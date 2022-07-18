package main

import (
	aio "iot_server6/aliyun_iot_client"
	pc "iot_server6/pulsar_client"
	"os"
	"os/signal"
	"syscall"

	"iot_server6/config"
	"iot_server6/log"
)

func main() {
	config.SetUp()
	log.InitLog()
	pc.SetUp()
	aio.BeginSubscribe()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch
	pc.Close()
}
