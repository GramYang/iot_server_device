package pulsar_client

import (
	"context"
	"encoding/json"
	"iot_server4/router"

	"github.com/apache/pulsar-client-go/pulsar"

	g "github.com/GramYang/gylog"
)

type aliyunResponse struct {
	Name    string `json:"name"`
	Product string `json:"product"`
	Time    string `json:"time"`
	Data    []int  `json:"data"`
}

func BeginSubscribe() {
	//处理iot_server6返回的阿里云响应
	go startReceiveMessage(consumer1, func(message pulsar.Message, c pulsar.Consumer) {
		var res aliyunResponse
		err := json.Unmarshal(message.Payload(), &res)
		if err != nil {
			g.Errorln(err)
			return
		}
		c.Ack(message)
		//目前只应用寄存器表1
		router.Router1(res.Name, res.Product, res.Time, res.Data)
	})
}

func startReceiveMessage(consumer pulsar.Consumer, handler func(pulsar.Message, pulsar.Consumer)) {
	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			g.Debugf("pulsar consumer %s receive data error%s\n", consumer.Name(), err)
		} else {
			handler(msg, consumer)
		}
	}
}
