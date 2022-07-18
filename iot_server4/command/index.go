package command

import (
	"encoding/json"
	"iot_server4/config"

	g "github.com/GramYang/gylog"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/iot"
)

func deviceWrite(name, productKey string, cmd []uint16) {
	client, err := iot.NewClientWithAccessKey(config.Conf.IotEndpoint, config.Conf.IotAccessKeyId, config.Conf.IotAccessKeySecret)
	if err != nil {
		g.Errorln(err)
		return
	}
	var modbus = make(map[string][]uint16)
	modbus["out"] = cmd
	data, err := json.Marshal(&modbus)
	request := iot.CreateSetDevicePropertyRequest()
	request.ProductKey = productKey
	request.DeviceName = name
	request.IotInstanceId = config.Conf.IotInstanceId
	request.Items = string(data)
	response, err := client.SetDeviceProperty(request)
	if err != nil {
		g.Errorln(err)
	} else {
		g.Debugf("响应 请求id %s 是否成功 %t\n 信息 %s", response.RequestId, response.Success, response.ErrorMessage)
	}
}
