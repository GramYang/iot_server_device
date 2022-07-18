package pulsar_client

import (
	"fmt"
	"iot_server4/config"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

var client pulsar.Client

//获取iot_server6分发的设备响应
var consumer1 pulsar.Consumer

func SetUp() {
	addr := fmt.Sprintf("pulsar://%s:%d", config.Conf.PulsarIp, config.Conf.PulsarPort)
	var err error
	client, err = pulsar.NewClient(pulsar.ClientOptions{URL: addr, ConnectionTimeout: 30 * time.Second})
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	consumer1, err = client.Subscribe(pulsar.ConsumerOptions{
		Topic: config.Conf.PulsarTopic1, SubscriptionName: "consumer1", Type: pulsar.Shared,
	})
	if err != nil {
		panic(err)
	}
}

func Close() {
	consumer1.Unsubscribe()
	consumer1.Close()
	client.Close()
}
