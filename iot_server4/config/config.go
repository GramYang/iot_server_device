package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"time"
)

var Conf = &Config{}
var configPath = "iot_server4_config"

type Config struct {
	ServerPort         int           `json:"server_port"`
	PromPort           int           `json:"prom_port"`
	LogLocal           bool          `json:"log_local"`
	LogDebug           bool          `json:"log_debug"`
	PulsarIp           string        `json:"pulsar_ip"`
	PulsarPort         int           `json:"pulsar_port"`
	PulsarTopic1       string        `json:"pulsar_topic1"` //接受iot_server6分发的设备消息
	JwtSecret          string        `json:"jwt_secret"`
	MysqlUserName      string        `json:"mysql_username"`
	MysqlPassword      string        `json:"mysql_password"`
	MysqlAddr          string        `json:"mysql_addr"`
	MysqlPort          int           `json:"mysql_port"`
	MysqlDatabase      string        `json:"mysql_database"`
	IotEndpoint        string        `json:"iot_endpoint"`
	IotAccessKeyId     string        `json:"iot_accesskeyid"`
	IotAccessKeySecret string        `json:"iot_accesskeysecret"`
	IotInstanceId      string        `json:"iot_instanceid"`
	CmdInterval        time.Duration `json:"cmd_interval"`       //同一设备命令发送间隙，推荐200毫秒
	LoopInterval       int           `json:"loop_interval"`      //主机和开关的循环指令的上传间隙，推荐4秒
	HeartbeatInterval  time.Duration `json:"heartbeat_interval"` //心跳时间，客户端和设备的心跳保持一致，单位秒
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
