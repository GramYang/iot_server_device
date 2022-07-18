package modbus

import (
	"time"
)

//主机下发命令————电箱初始化
func HostWrite1() []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, 0, HOST_BASEADDR+HOSTREGISTER_ISSUEDCMD, HOSTREGISTER_ISSUEDCMD_INIT)
	return tmp
}

//主机下发命令————电箱变动
func HostWrite2() []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, 0, HOST_BASEADDR+HOSTREGISTER_ISSUEDCMD, HOSTREGISTER_ISSUEDCMD_CHANGE)
	return tmp
}

//主机定时上报间隔，这个上报间隔单位为秒，默认值5分钟，过了5分钟会失效，需要你每5分钟就设置一次值才能保持
func HostWrite3(data uint16) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, 0, HOST_BASEADDR+HOSTREGISTER_UPLOADINTERVAL, data)
	return tmp
}

//主机下发命令————电箱地址变动
func HostWrite4() []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, 0, HOST_BASEADDR+HOSTREGISTER_ISSUEDCMD, HOSTREGISTER_ISSUEDCMD_CHANGEADDR)
	return tmp
}

//一键开合闸命令加密，这里只是打印一下报文
func HostAllSwitch(stateArr []int) []uint16 {
	tmp := allSwitchEncode(stateArr)
	return tmp
}

//返回时间戳
func HostTimestamp() []uint16 {
	tmp := make([]byte, 5)
	tmp[0] = MSGTYPE_TIMESTAMP_RESULT
	putUint32Switch(tmp[1:5], uint32(time.Now().Unix()), BIGEND)
	tmp = append(tmp)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(tmp), CRC_BIGEND)
	tmp = append(tmp, crc...)
	return byte2uint16(tmp)
}

//设备命令——合闸，开启
func DeviceCMD1(addr uint8) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+DEVICEREGISTER_ONE_CMD, DEVICEREGISTER_ONE_CMD_REMOTEOPEN)
	return tmp
}

//设备命令——关闭
func DeviceCMD2(addr uint8) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+DEVICEREGISTER_ONE_CMD, DEVICEREGISTER_ONE_CMD_REMOTECLOSE)
	return tmp
}

//设备命令——锁定
func DeviceCMD3(addr uint8) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+DEVICEREGISTER_ONE_CMD, DEVICEREGISTER_ONE_CMD_REMOTELOCK)
	return tmp
}

//设备命令——解锁
func DeviceCMD4(addr uint8) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+DEVICEREGISTER_ONE_CMD, DEVICEREGISTER_ONE_CMD_REMOTEUNLOCK)
	return tmp
}

//设备命令——漏电测试
func DeviceCMD5(addr uint8) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+DEVICEREGISTER_ONE_CMD, DEVICEREGISTER_ONE_CMD_ELECTRICLEAKAGETEST)
	return tmp
}

//设备命令——当前故障清除
func DeviceCMD6(addr uint8) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+DEVICEREGISTER_ONE_CMD, DEVICEREGISTER_ONE_CMD_CLEARFAULT)
	return tmp
}

//设备命令——阈值恢复出厂设置
func DeviceCMD7(addr uint8) []uint16 {
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, DEVICE_BASEADDR+DEVICEREGISTER_ONE_CMD, DEVICEREGISTER_ONE_CMD_RESET)
	return tmp
}

//设备控制——预警开启——全功能(没有开启的全部都关闭)
func DeviceControlPreAlarmTotal(addr uint8, registerTable uint8, enableList []uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ALARMENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ALARMENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitsToData(enableList))
	return tmp
}

//设备控制——预警开启——过流
func DeviceControlPreAlarm1(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ALARMENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ALARMENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ALARMENABLE_OVERFLOW))
	return tmp
}

//设备控制——预警开启——过压
func DeviceControlPreAlarm2(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ALARMENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ALARMENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ALARMENABLE_OVERVOLTAGE))
	return tmp
}

//设备控制——预警开启——欠压
func DeviceControlPreAlarm3(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ALARMENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ALARMENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ALARMENABLE_UNDERVOLTAGE))
	return tmp
}

//设备控制——预警开启——过载
func DeviceControlPreAlarm5(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ALARMENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ALARMENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ALARMENABLE_OVERDRIVE))
	return tmp
}

//设备控制——预警开启——电量
func DeviceControlPreAlarm6(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ALARMENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ALARMENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ALARMENABLE_ELECTRICQUANTITY))
	return tmp
}

//设备控制——预警开启——过温
func DeviceControlPreAlarm7(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ALARMENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ALARMENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ALARMENABLE_OVERHEAT))
	return tmp
}

//设备控制——故障开启——全功能(没有开启的全部都关闭)
func DeviceControlFaultProtectTotal(addr uint8, registerTable uint8, enableList []uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitsToData(enableList))
	return tmp
}

//设备控制——故障保护开启——过流
func DeviceControlFaultProtect1(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_OVERFLOW))
	return tmp
}

//设备控制——故障保护开启——过压
func DeviceControlFaultProtect2(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_OVERVOLTAGE))
	return tmp
}

//设备控制——故障保护开启——欠压
func DeviceControlFaultProtect3(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_UNDERVOLTAGE))
	return tmp
}

//设备控制——故障保护开启——过载
func DeviceControlFaultProtect5(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_OVERDRIVE))
	return tmp
}

//设备控制——故障保护开启——电量
func DeviceControlFaultProtect6(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_ELECTRICQUANTITY))
	return tmp
}

//设备控制——故障保护开启——过温
func DeviceControlFaultProtect7(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_OVERHEAT))
	return tmp
}

//设备控制——故障保护开启——短路
func DeviceControlFaultProtect9(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_SHORTOUT))
	return tmp
}

//设备控制——故障保护开启——电弧
func DeviceControlFaultProtect12(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_ERRORENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_ERRORENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, bitToData(DEVICEREGISTER_ONE_ERRORENABLE_ELECTRICARC))
	return tmp
}

//设备控制——过欠压恢复开启——关
func DeviceControlVolLimitRst1(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_VOLTAGELIMITRSTENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_VOLTAGELIMITRSTENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, DEVICEREGISTER_ONE_VOLTAGELIMITRSTENABLE_CLOSE)
	return tmp
}

//设备控制——过欠压恢复开启——开
func DeviceControlVolLimitRst2(addr uint8, registerTable uint8) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_VOLTAGELIMITRSTENABLE
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_VOLTAGELIMITRSTENABLE
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, DEVICEREGISTER_ONE_VOLTAGELIMITRSTENABLE_OPEN)
	return tmp
}

//设备控制——限定电流预警值设置
func DeviceControl1(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_IH_P
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_IH_P
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——限定电流保护值设置
func DeviceControl2(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_IH
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_IH
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——过压预警值设置
func DeviceControl3(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_UH_P
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_UH_P
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——过压保护值设置
func DeviceControl4(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_UH
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_UH
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——欠压预警值设置
func DeviceControl5(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_UL_P
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_UL_P
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——欠压保护值设置
func DeviceControl6(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_UL
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_UL
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——漏电流预警值设置
func DeviceControl9(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_IL_P
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_IL_P
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——漏电流保护值设置
func DeviceControl10(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_IL
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_IL
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——限定功率预警值设置
func DeviceControl11(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_PH_P
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_PH_P
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——限定功率保护值设置
func DeviceControl12(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_PH
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_PH
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——剩余电量预警值设置
func DeviceControl13(addr uint8, registerTable uint8, value uint32) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_EH_P
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_EH_P
	}
	tmp := writeSingle32Encode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——冲减电量
func DeviceControl14(addr uint8, registerTable uint8, index uint8, value int) []uint16 {
	data := make([]byte, 4)
	var realValue uint16
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_EH
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_EH
	}
	if value > 0 {
		data[3] = 0xa5
		realValue = uint16(value)
	} else if value < 0 {
		data[3] = 0x5a
		realValue = uint16(-value)
	} else {
		return nil
	}
	data[2] = index
	putUint16Switch(data[0:2], realValue, BIGEND)
	tmp := writeSingle32Encode(MSGTYPE_CMD, addr, registerNo, uint32Switch(data, BIGEND))
	return tmp
}

//设备控制——过温预警值设置，必须是35-125
func DeviceControl15(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_TH_P
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_TH_P
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——过温保护值设置，必须是35-125
func DeviceControl16(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_TH
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_TH
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——过欠压动作时间值设置，必须是0-60
func DeviceControl19(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_UHL_CT
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_UHL_CT
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——过欠压恢复时间值设置，必须是10-60
func DeviceControl20(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_UHL_RT
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_UHL_RT
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备控制——电流功率动作时间（秒，0-60）
func DeviceControl21(addr uint8, registerTable uint8, value uint16) []uint16 {
	var registerNo uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_IH_PH_CT
	case REGISTER_TABLE_TWO:
		registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_IH_PH_CT
	}
	tmp := writeSingleEncode(MSGTYPE_CMD, addr, registerNo, value)
	return tmp
}

//设备定时器设置——单次
func DeviceTimerOnce(addr, registerTable, group, timerType, state uint8, timestamp uint32) []uint16 {
	tmp := deviceTimerOnce(addr, registerTable, group, timerType, state, timestamp)
	return tmp
}

//设备定时器设置——每天
func DeviceTimerDaily(addr, registerTable, group, timerType, state, hour, minute uint8) []uint16 {
	tmp := deviceTimerDaily(addr, registerTable, group, timerType, state, hour, minute)
	return tmp
}

//设备定时器设置——每周
func DeviceTimerWeekly(addr, registerTable, group, timerType, state, weekday, hour, minute uint8) []uint16 {
	tmp := deviceTimerWeekly(addr, registerTable, group, timerType, state, weekday, hour, minute)
	return tmp
}

//设备定时器设置——每月
func DeviceTimerMonthly(addr, registerTable, group, timerType, state, day, hour, minute uint8) []uint16 {
	tmp := deviceTimerMonthly(addr, registerTable, group, timerType, state, day, hour, minute)
	return tmp
}
