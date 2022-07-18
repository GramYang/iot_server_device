package command

import (
	"iot_server4/modbus"
)

//主机地址重分配并锁定设备
func HostResetAndLock(deviceid, product string) {
	deviceWrite(deviceid, product, modbus.HostWrite2())
}

//主机地址重分配不锁定设备
func HostReset(deviceid, product string) {
	deviceWrite(deviceid, product, modbus.HostWrite4())
}

//设置主机快照上传周期
func HostUploadInterval(deviceid, product string, second int) {
	deviceWrite(deviceid, product, modbus.HostWrite3(uint16(second)))
}

//主机一键分合闸
func HostSwitch(deviceid, product string, cmd []int) {
	deviceWrite(deviceid, product, modbus.HostAllSwitch(cmd))
}

//读取主机快照（发送间隔不能低于3秒，发送次数0表示关停发送，两者均是无符号整数）
//响应只有设备的有用功电压和电量
func HostSnapShoot(deviceid, product string, interval, num uint8) {
	deviceWrite(deviceid, product, modbus.ReadSnapShoot(interval, num))
}

//读主机全寄存器（6-10）
//响应有版本号、信号强度、通信方式、子设备数、主机上传间隔
func HostAllRegister(deviceid, product string) {
	deviceWrite(deviceid, product, modbus.ReadHost1())
}

//读iccid
func Iccid(deviceid, product string) {
	deviceWrite(deviceid, product, modbus.HostIccid())
}

//读主机下面所有设备状态
//设备状态：在线、硬件故障、预警状态、故障状态、开关状态、开关模式
func HostAllSwitchState(deviceid, product string) {
	deviceWrite(deviceid, product, modbus.ReadHostAllState())
}

//响应主机时间戳请求
func HostTimestampResponse(deviceid, product string) {
	deviceWrite(deviceid, product, modbus.HostTimestamp())
}
