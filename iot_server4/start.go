package main

import (
	"iot_server4/cache"
	"iot_server4/config"
	"iot_server4/handler"
	"iot_server4/heartbeat"
	"iot_server4/log"
	"iot_server4/model"
	"iot_server4/prom"
	pc "iot_server4/pulsar_client"
	q "iot_server4/queue"
	sc "iot_server4/sqlx_client"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"

	_ "iot_server4/proc"

	_ "github.com/davyxu/cellnet/peer/gorillaws"
)

func main() {
	config.SetUp()
	cache.SetUp()
	log.InitLog()
	pc.SetUp()
	sc.SetUp()
	pc.BeginSubscribe()
	prom.SetUp()
	heartbeat.StartCheck()
	q.SetUp(495)
	address := "ws://:" + strconv.Itoa(config.Conf.ServerPort) + "/entry"
	p := peer.NewGenericPeer("gorillaws.Acceptor", "iot_server4", address, nil)
	proc.BindProcessorHandler(p, "iot_server4", handler.Handler2)
	if socketOpt, ok := p.(cellnet.TCPSocketOption); ok {
		socketOpt.SetSocketBuffer(2048, 2048, true)
		socketOpt.SetSocketDeadline(time.Second*40, time.Second*20)
	}
	p.Start()
	model.FrontendSessionManager = p.(peer.SessionManager)
	model.AddLocalService(p)
	model.CheckReady()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-ch
	model.StopAllService()
	pc.Close()
}
