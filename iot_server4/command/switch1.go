package command

import "iot_server4/modbus"

//开关开启
func SwitchOn(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.DeviceCMD1(addr))
}

//开关关闭
func SwitchOff(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.DeviceCMD2(addr))
}

//开关锁定
func SwitchLockOn(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.DeviceCMD3(addr))
}

//开关解锁
func SwitchLockOff(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.DeviceCMD4(addr))
}

//开关漏电测试
func SwitchElectricLeakageTest(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.DeviceCMD5(addr))
}

//读开关实时数据（发送间隔不能低于3秒，发送次数0表示关停发送，两者均是无符号整数）
func SwitchRuntime(deviceId, product string, addr, interval, count uint8) {
	deviceWrite(deviceId, product, modbus.ReadDevice2(addr, interval, count))
}

//读开关实时数据（modbus范围读，单次）
func SwitchRuntimeOnce(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.ReadDevice5(addr))
}

//读开关设置数据
func SwitchSetting(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.ReadDevice3(addr))
}

//开关当前故障清除
func SwitchClearCurrentError(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.DeviceCMD6(addr))
}

//开关阈值恢复出厂设置
func SwitchResetThresholdValue(deviceId, product string, addr uint8) {
	deviceWrite(deviceId, product, modbus.DeviceCMD7(addr))
}

//预警开启——过流
func SwitchAlarmEnable0(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlPreAlarm1(addr, 1))
}

//预警开启——过压
func SwitchAlarmEnable1(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlPreAlarm2(addr, 1))
}

//预警开启——欠压
func SwitchAlarmEnable2(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlPreAlarm3(addr, 1))
}

//预警开启——过载
func SwitchAlarmEnable3(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlPreAlarm5(addr, 1))
}

//预警开启——电量
func SwitchAlarmEnable4(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlPreAlarm6(addr, 1))
}

//预警开启——过温
func SwitchAlarmEnable5(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlPreAlarm7(addr, 1))
}

//预警开启——全
func SwitchAlarmEnableTotal(deviceid, product string, addr uint8, list []uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlPreAlarmTotal(addr, 1, list))
}

//故障保护开启——过流
func SwitchErrorEnable0(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtect1(addr, 1))
}

//故障保护开启——过压
func SwitchErrorEnable1(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtect2(addr, 1))
}

//故障保护开启——欠压
func SwitchErrorEnable2(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtect3(addr, 1))
}

//故障保护开启——过载
func SwitchErrorEnable3(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtect5(addr, 1))
}

//故障保护开启——电量
func SwitchErrorEnable4(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtect6(addr, 1))
}

//故障保护开启——过温
func SwitchErrorEnable5(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtect7(addr, 1))
}

//故障保护开启——电弧
func SwitchErrorEnable6(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtect12(addr, 1))
}

//故障保护开启——全
func SwitchErrorEnableTotal(deviceid, product string, addr uint8, list []uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlFaultProtectTotal(addr, 1, list))
}

//过压欠压恢复开启——开
func SwitchVolLimitRst0(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlVolLimitRst2(addr, 1))
}

//过压欠压恢复开启——关
func SwitchVolLimitRst1(deviceid, product string, addr uint8) {
	deviceWrite(deviceid, product, modbus.DeviceControlVolLimitRst1(addr, 1))
}

//限定电流预警值
func Switch_IH_P(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl1(addr, 1, value))
}

//限定电流保护值
func Switch_IH(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl2(addr, 1, value))
}

//过压预警值
func Switch_UH_P(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl3(addr, 1, value))
}

//过压保护值，低于250无效
func Switch_UH(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl4(addr, 1, value))
}

//欠压预警值，低于187无效
func Switch_UL_P(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl5(addr, 1, value))
}

//欠压保护值，低于187无效
func Switch_UL(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl6(addr, 1, value))
}

//漏电流预警值
func Switch_IL_P(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl9(addr, 1, value))
}

//漏电流保护值
func Switch_IL(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl10(addr, 1, value))
}

//限定功率预警值
func Switch_PH_P(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl11(addr, 1, value))
}

//限定功率保护值，低于4400无效
func Switch_PH(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl12(addr, 1, value))
}

//剩余电量预警值
func Switch_EH_P(deviceid, product string, addr uint8, value uint32) {
	deviceWrite(deviceid, product, modbus.DeviceControl13(addr, 1, value))
}

//冲减电量，这个接口目前有问题，先不管
func Switch_EH(deviceid, product string, addr uint8, index uint8, value int) {
	deviceWrite(deviceid, product, modbus.DeviceControl14(addr, 1, index, value))
}

//过温预警值
func Switch_TH_P(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl15(addr, 1, value))
}

//过温保护值，低于50就会重置成85
func Switch_TH(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl16(addr, 1, value))
}

//过欠压动作时间
func Switch_UHL_CT(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl19(addr, 1, value))
}

//过欠压恢复时间
func Switch_UHL_RT(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl20(addr, 1, value))
}

//电流功率动作时间
func Switch_IH_PH_CT(deviceid, product string, addr uint8, value uint16) {
	deviceWrite(deviceid, product, modbus.DeviceControl21(addr, 1, value))
}

//设置定时器——单次
func Switch_SetTimer_Once(deviceid, product string, addr uint8, group, timerType, state uint8, timestamp uint32) {
	deviceWrite(deviceid, product, modbus.DeviceTimerOnce(addr, 1, group, timerType, state, timestamp))
}

//设置定时器——每天
func Switch_SetTimer_Daily(deviceid, product string, addr uint8, group, timerType, state, hour, minute uint8) {
	deviceWrite(deviceid, product, modbus.DeviceTimerDaily(addr, 1, group, timerType, state, hour, minute))
}

//设置定时器——每周
func Switch_SetTimer_Weekly(deviceid, product string, addr uint8, group, timerType, state, weekday, hour, minute uint8) {
	deviceWrite(deviceid, product, modbus.DeviceTimerWeekly(addr, 1, group, timerType, state, weekday, hour, minute))
}

//设置定时器——每月
func Switch_SetTimer_Monthly(deviceid, product string, addr uint8, group, timerType, state, day, hour, minute uint8) {
	deviceWrite(deviceid, product, modbus.DeviceTimerMonthly(addr, 1, group, timerType, state, day, hour, minute))
}
