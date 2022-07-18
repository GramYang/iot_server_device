package modbus

import (
	"fmt"
	"strings"
)

//读主机全寄存器（6-10）
func ReadHost1() []uint16 {
	tmp := readRangeEncode(MSGTYPE_CMD, 0, HOST_BASEADDR+6, 5)
	return tmp
}

//主机快照读全部寄存器
func ReadHostSnap1() []uint16 {
	tmp := snapShootReadRangeEncode(1, 1)
	return tmp
}

//一键分合闸响应读取
func ReadHostAllSwitch(data []uint16) ([]int, error) {
	return allSwitchDecode(data)
}

//读设备寄存器——实时数据
func ReadDevice2(addr, interval, count uint8) []uint16 {
	tmp := []uint16{uint16(MSGTYPE_RUNTIME), uint16(addr), uint16(interval), uint16(count)}
	return tmp
}

//读设备寄存器——设置字段
func ReadDevice3(addr uint8) []uint16 {
	tmp := []uint16{uint16(MSGTYPE_SETTING), uint16(addr)}
	return tmp
}

//读设备寄存器——定时器状态
func ReadDevice4(addr uint8) []uint16 {
	tmp := readRangeEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+47, 32)
	return tmp
}

//读设备寄存器——实时数据
func ReadDevice5(addr uint8) []uint16 {
	tmp := readRangeEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+3, 24)
	return tmp
}

//设备实时数据定时上报解析
func ReadDeviceRuntimeUploadResult(data []uint16) *DeviceStatusUpload {
	result := &DeviceStatusUpload{}
	result.LeftTime = uint8(data[1])
	result.Addr = uint8(data[2])
	result.DeviceType = uint8(data[3])
	result.RatedCurrent = uint16Switch(uint162byte(data[4:6]), BIGEND)
	if DEVICETYPE_REGISTERTYPE[result.DeviceType] == REGISTER_TABLE_ONE {
		t1 := &DeviceRegisterDataRead1{}
		intSliceMapBean(data[6:54], t1)
		result.Data = *t1
	} else if DEVICETYPE_REGISTERTYPE[result.DeviceType] == REGISTER_TABLE_TWO {
		t2 := &DeviceRegisterDataRead2{}
		intSliceMapBean(data[6:58], t2)
		result.Data = *t2
	}
	return result
}

//主机快照定时上报解析
func ReadSnapShootUploadResult(data []uint16) *SnapShootReadRangeResult {
	return SnapShootTimedUploadDecode(data[1:])
}

//主机快照读取
func ReadSnapShoot(interval, mount uint8) []uint16 {
	tmp := []uint16{uint16(MSGTYPE_GETSNAPSHOOT), uint16(interval), uint16(mount)}
	return tmp
}

//主机快照读取解析
func ReadSnapShootResult(data []uint16) *SnapShootResult {
	return SnapShootMultiUploadDecode(data[1:])
}

//设备上传设置字段解析
func ReadDeviceSettingUploadResult(data []uint16) *DeviceSettingStatus {
	result := &DeviceSettingStatus{}
	result.Addr = uint8(data[1])
	result.DeviceType = uint8(data[2])
	result.RatedCurrent = uint16Switch(uint162byte(data[3:5]), BIGEND)
	t1 := &DeviceRegisterDataWrite1{}
	intSliceMapBean(data[5:49], t1)
	result.Data = *t1
	var t2 []DeviceTimer1
	for i := 49; i < 113; i += 16 {
		t2 = append(t2, *deviceTimerDecode1(data[i : i+16]))
	}
	result.Timers = t2
	return result
}

//故障事件上报解析
func DeviceFaultEventUploadResult(data []uint16) *DeviceFaultEvent {
	result := &DeviceFaultEvent{}
	result.Addr = uint8(data[1])
	result.DeviceType = uint8(data[2])
	result.RatedCurrent = uint16Switch(uint162byte(data[3:5]), BIGEND)
	result.FaultEvents = uint32Switch(uint162byte(data[5:9]), BIGEND)
	t1 := &FaultData1{}
	intSliceMapBean(data[9:37], t1)
	result.Data = *t1
	return result
}

//预警事件上报解析
func DeviceWarnEventUploadResult(data []uint16) *DeviceWarnEvent {
	result := &DeviceWarnEvent{}
	result.Addr = uint8(data[1])
	result.DeviceType = uint8(data[2])
	result.RatedCurrent = uint16Switch(uint162byte(data[3:5]), BIGEND)
	result.WarnEvents = uint32(uint16Switch(uint162byte(data[5:9]), BIGEND))
	t1 := &WarnData1{}
	intSliceMapBean(data[9:37], t1)
	result.Data = *t1
	return result
}

//硬件故障事件上报解析
func DeviceHardwareFaultUploadResult(data []uint16) *DeviceHardwareFaultEvent {
	result := &DeviceHardwareFaultEvent{}
	result.Addr = uint8(data[1])
	result.DeviceType = uint8(data[2])
	result.RatedCurrent = uint16Switch(uint162byte(data[3:5]), BIGEND)
	result.HardwareEvents = uint16Switch(uint162byte(data[5:7]), BIGEND)
	return result
}

//开关状态事件上报解析
func SwitchStatusEventUploadResult(data []uint16) *SwitchStatusEvent {
	result := &SwitchStatusEvent{}
	result.Addr = uint8(data[1])
	result.DeviceType = uint8(data[2])
	result.RatedCurrent = uint16Switch(uint162byte(data[3:5]), BIGEND)
	result.SwitchStatus = uint16Switch(uint162byte(data[5:7]), BIGEND)
	return result
}

//开关模式事件上报解析
func SwitchModeEventUploadResult(data []uint16) *SwitchModeEvent {
	result := &SwitchModeEvent{}
	result.Addr = uint8(data[1])
	result.DeviceType = uint8(data[2])
	result.RatedCurrent = uint16Switch(uint162byte(data[3:5]), BIGEND)
	result.SwitchMode = uint16Switch(uint162byte(data[5:7]), BIGEND)
	return result
}

//漏电检测事件上报解析
func ElectricLeakageTestEventUploadResult(data []uint16) *ElectricLeakageTestEvent {
	result := &ElectricLeakageTestEvent{}
	result.Addr = uint8(data[1])
	result.DeviceType = uint8(data[2])
	result.RatedCurrent = uint16Switch(uint162byte(data[3:5]), BIGEND)
	switch data[5] {
	case 0:
		result.Status = "有漏电"
	case 1:
		result.Status = "无漏电"
	}
	return result
}

//主机信息事件上报解析
func HostInfoEventUploadResult(data []uint16) *HostRegisterData {
	p := &HostRegisterData{}
	intSliceMapBean(data[:10], p)
	return p
}

//读主机iccid
func HostIccid() []uint16 {
	tmp := []uint16{uint16(MSGTYPE_ICCID)}
	return tmp
}

//读主机iccid响应解析
func HostIccidResult(data []uint16) string {
	var sb strings.Builder
	for _, v := range data {
		sb.WriteString(fmt.Sprintf("%c", v))
	}
	return sb.String()
}

//读主机下所有设备状态信息
func ReadHostAllState() []uint16 {
	tmp := []uint16{uint16(MSGTYPE_ALLSTATE)}
	return tmp
}

//读主机下所有设备状态信息，响应解析
func ReadHostAllStateResult(data []uint16) *SnapShootReadRangeResult {
	result := &SnapShootReadRangeResult{}
	result.Ok = true
	result.DeviceCount = uint8(data[0])
	var index = 1
	for i := 0; i < int(data[0]); i++ {
		snapShoot := DeviceSnapShoot{}
		snapShoot.DeviceAddr = uint8(data[index])
		snapShoot.DeviceType = uint8(data[index+1])
		snapShoot.RatedCurrent = uint16Switch(uint162byte(data[index+2:index+4]), BIGEND)
		state := DeviceState{}
		intSliceMapBean(data[index+4:index+19], &state)
		snapShoot.Data = state
		result.DeviceSnapShoots = append(result.DeviceSnapShoots, snapShoot)
		index += 19
	}
	return result
}

//电量操作事件解析
func ReadDeviceElectricQuantityChangeEvent(data []uint16) {
	var resultMsg string
	switch data[4] {
	case 1:
		resultMsg = "充电成功"
	case 2:
		resultMsg = "充电失败"
	case 3:
		resultMsg = "扣电成功"
	case 4:
		resultMsg = "扣电失败"
	case 5:
		resultMsg = "操作命令错误"
	}
	fmt.Printf("设备地址:%d 设备类型:%s 额定电流:%d 操作结果:%s 序列号:%d 电量值:%d\n",
		data[0], DEVICETYPE_MSG[uint8(data[1])], uint16Switch(uint162byte(data[2:4]), BIGEND), resultMsg, data[5], uint16Switch(uint162byte(data[6:8]), BIGEND))
}
