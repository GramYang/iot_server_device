package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var Conf = &Config{}
var configPath = "iot_server6_config"

type Config struct {
	ServerPort           int    `json:"server_port"`
	LogLocal             bool   `json:"log_local"`
	LogDebug             bool   `json:"log_debug"`
	PulsarIp             string `json:"pulsar_ip"`
	PulsarPort           int    `json:"pulsar_port"`
	PulsarTopic          string `json:"pulsar_topic"`
	IotEndpoint          string `json:"iot_endpoint"`
	IotAccessKeyId       string `json:"iot_accesskeyid"`
	IotAccessKeySecret   string `json:"iot_accesskeysecret"`
	IotInstanceId        string `json:"iot_instanceid"`
	IotConsumerGroupId   string `json:"iot_consumer_groupid"`
	IotUid               string `json:"iot_uid"`
	IotSubscribeClientId string `json:"iot_subscribe_clientid"`
}

func SetUp() {
	var p string
	flag.StringVar(&p, "c", "", "配置文件路径")
	flag.Parse()
	if p == "" {
		p = configPath
	}
	file, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(file, Conf); err != nil {
		panic(err)
	}
}
