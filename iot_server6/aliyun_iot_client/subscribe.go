package aliyun_iot_client

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"pack.ag/amqp"

	g "github.com/GramYang/gylog"

	"iot_server6/config"
	"iot_server6/router"
)

var timestampFilter int64

type AmqpManager struct {
	address  string
	userName string
	password string
	client   *amqp.Client
	session  *amqp.Session
	receiver *amqp.Receiver
}

type MessageDecode struct {
	DeviceType      string      `json:"deviceType"`
	IotId           string      `json:"iotId"`
	RequestId       string      `json:"requestId"`
	CheckFailedData interface{} `json:"checkFailedData"`
	ProductKey      string      `json:"productKey"`
	GmtCreate       int64       `json:"gmtCreate"`
	DeviceName      string      `json:"deviceName"`
	Items           `json:"items"`
}

type Items struct {
	In `json:"in"`
}

type In struct {
	Value []int `json:"value"`
	Time  int64 `json:"time"`
}

func (am *AmqpManager) generateReceiver() error {
	//如果session不为空则重新连接
	if am.session != nil {
		receiver, err := am.session.NewReceiver(
			amqp.LinkSourceAddress("/queue1"),
			amqp.LinkCredit(20),
		)
		if err != nil {
			return err
		} else {
			am.receiver = receiver
			return nil
		}
	}
	//client不为空则清理掉
	if am.client != nil {
		am.client.Close()
	}
	//重新连接
	client, err := amqp.Dial(am.address, amqp.ConnSASLPlain(am.userName, am.password))
	if err != nil {
		return err
	}
	am.client = client
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	am.session = session
	receiver, err := am.session.NewReceiver(
		amqp.LinkSourceAddress("/queue1"),
		amqp.LinkCredit(20),
	)
	if err != nil {
		return err
	}
	am.receiver = receiver
	return nil
}

//连接重试
func (am *AmqpManager) generateReceiverWithRetry(ctx context.Context) error {
	duration := 10 * time.Millisecond
	maxDuration := 20 * time.Second
	times := 1
	for {
		err := am.generateReceiver()
		if err != nil {
			time.Sleep(duration)
			if duration < maxDuration {
				duration *= 2
			}
			g.Debugf("amqp connect retry, times:%d, duration:%d\n", times, duration)
			times++
		} else {
			g.Debugln("amqp connect init success")
			return nil
		}
	}
}

func sliceFmt(data []int) string {
	var b strings.Builder
	b.WriteString("[")
	for k, v := range data {
		b.WriteString(strconv.Itoa(v))
		if k != len(data)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("]")
	return b.String()
}

func (am *AmqpManager) processMessage(message *amqp.Message) {
	g.Debugf("%#v\n", message)
	g.Debugf("data received:%s properties:%v\n", string(message.GetData()), message.ApplicationProperties)
	m := &MessageDecode{}
	_ = json.Unmarshal(message.GetData(), m)
	//过滤一下乱序分发的老消息
	if m.Time < timestampFilter {
		return
	}
	receivedTimestamp := time.Unix(m.Time/1000, 0).String()
	g.Debugf("device: %s timestamp: %s value: %s\n", m.DeviceName, receivedTimestamp, sliceFmt(m.Value))
	router.Router(m.DeviceName, m.ProductKey, receivedTimestamp, m.Value)
}

func (am *AmqpManager) startReceiveMessage() {
	ctx := context.Background()
	err := am.generateReceiverWithRetry(ctx)
	if err != nil {
		return
	}
	defer func() {
		am.receiver.Close(ctx)
		am.session.Close(ctx)
		am.client.Close()
	}()
	for {
		message, err := am.receiver.Receive(ctx)
		if err == nil {
			go am.processMessage(message)
			message.Accept()
		} else {
			g.Debugln("amqp receive data error:", err)
			err := am.generateReceiverWithRetry(ctx)
			if err != nil {
				return
			}
		}
	}
}

func BeginSubscribe() {
	host := fmt.Sprintf("%s.iot-amqp.%s.aliyuncs.com", config.Conf.IotUid, config.Conf.IotEndpoint)
	address := "amqps://" + host + ":5671"
	timestamp := time.Now().UnixMilli()
	userName := fmt.Sprintf("%s|authMode=aksign,signMethod=Hmacsha1,consumerGroupId=%s,authId=%s,iotInstanceId=%s,timestamp=%d|",
		config.Conf.IotSubscribeClientId, config.Conf.IotConsumerGroupId, config.Conf.IotAccessKeyId, config.Conf.IotInstanceId, timestamp)
	stringToSign := fmt.Sprintf("authId=%s&timestamp=%d", config.Conf.IotAccessKeyId, timestamp)
	hmacKey := hmac.New(sha1.New, []byte(config.Conf.IotAccessKeySecret))
	hmacKey.Write([]byte(stringToSign))
	password := base64.StdEncoding.EncodeToString(hmacKey.Sum(nil))
	amqpManager := &AmqpManager{
		address:  address,
		userName: userName,
		password: password,
	}
	amqpManager.startReceiveMessage()
	timestampFilter = time.Now().UnixMilli()
}
