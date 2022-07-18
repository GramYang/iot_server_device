package main

import (
	"fmt"
	"iot_server4/modbus"
)

func main() {
	arr := []uint16{184, 0, 1, 1, 63, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 136, 1, 207, 21, 3, 0, 3, 166}
	res := modbus.ReadDeviceRuntimeUploadResult(arr)
	fmt.Printf("%#v\n", res)
}
