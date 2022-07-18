package pulsar_client

import (
	"context"
	"fmt"
	"iot_server6/config"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"

	g "github.com/GramYang/gylog"
)

var client pulsar.Client
var producer pulsar.Producer

func SetUp() {
	addr := fmt.Sprintf("pulsar://%s:%d", config.Conf.PulsarIp, config.Conf.PulsarPort)
	var err error
	client, err = pulsar.NewClient(pulsar.ClientOptions{URL: addr, ConnectionTimeout: 30 * time.Second})
	if err != nil {
		panic(err)
	}
	producer, err = client.CreateProducer(pulsar.ProducerOptions{
		Topic: config.Conf.PulsarTopic,
	})
	if err != nil {
		panic(err)
	}
}

func Close() {
	producer.Close()
	client.Close()
}

func Send(bs []byte) {
	_, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: bs,
	})
	if err != nil {
		g.Errorln(err)
	}
}
