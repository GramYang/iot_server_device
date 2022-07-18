package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var Conf = &Config{}
var configPath = "iot_server8_config"

type Config struct {
	ServerPort    int    `json:"server_port"`
	LogLocal      bool   `json:"log_local"`
	LogDebug      bool   `json:"log_debug"`
	MysqlUserName string `json:"mysql_username"`
	MysqlPassword string `json:"mysql_password"`
	MysqlAddr     string `json:"mysql_addr"`
	MysqlPort     int    `json:"mysql_port"`
	MysqlDatabase string `json:"mysql_database"`
	JwtSecret     string `json:"jwt_secret"`
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
