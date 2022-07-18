package router

import (
	"encoding/json"
	pc "iot_server6/pulsar_client"
)

func Router(deviceId, product, receivedTimestamp string, data []int) {
	content := map[string]interface{}{"name": deviceId, "product": product, "time": receivedTimestamp, "data": data}
	bs, _ := json.Marshal(content)
	pc.Send(bs)
}
