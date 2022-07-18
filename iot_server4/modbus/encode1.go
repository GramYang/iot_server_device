package modbus

import (
	"errors"
)

//设备&主机寄存器单独写加密
func writeSingleEncode(msgType uint8, addr uint8, registerNo uint16, data uint16) []uint16 {
	var result []uint16
	result = append(result, []uint16{uint16(msgType), uint16(addr), 6}...)
	src := []uint16{registerNo, data}
	dst := make([]byte, 4)
	for i := 0; i < 4; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	result = append(result, byte2uint16(dst)...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}

//设备&主机寄存器单独写加密，写大寄存器32位值，响应解析使用writeRangeDecode
func writeSingle32Encode(msgType uint8, addr uint8, registerStart uint16, data uint32) []uint16 {
	var result []uint16
	result = append(result, []uint16{uint16(msgType), uint16(addr), 16}...)
	src := []uint16{registerStart, 2}
	dst := make([]byte, 8)
	for i := 0; i < 4; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	putUint32Switch(dst[4:8], data, BIGEND)
	result = append(result, byte2uint16(dst[0:4])...)
	result = append(result, 4)
	result = append(result, byte2uint16(dst[4:8])...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}

//设备&主机寄存器范围写加密
func writeRangeEncode(msgType uint8, addr uint8, registerStart uint16, registerCount uint16, dataSize uint8, data []uint16) ([]uint16, error) {
	if len(data) != int(dataSize) {
		return nil, errors.New("dataSize wrong")
	}
	var result []uint16
	result = append(result, []uint16{uint16(msgType), uint16(addr), 16}...)
	src := []uint16{registerStart, registerCount}
	src = append(src, data...)
	dstLen := (len(data) + 2) * 2
	dst := make([]byte, dstLen)
	for i := 0; i < dstLen; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	for i := 0; i < dstLen; i += 2 {
		if i == 4 {
			result = append(result, uint16(dataSize))
		}
		result = append(result, byte2uint16(dst[i:i+2])...)
	}
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result, nil
}

//设备&主机寄存器范围读加密
func readRangeEncode(msgType uint8, addr uint8, registerStart uint16, registerCount uint16) []uint16 {
	var result []uint16
	result = append(result, []uint16{uint16(msgType), uint16(addr), 3}...)
	src := []uint16{registerStart, registerCount}
	dst := make([]byte, 4)
	for i := 0; i < 4; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	result = append(result, byte2uint16(dst)...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}

//主机快照范围读加密，没有消息类型
func snapShootReadRangeEncode(deviceStart uint8, deviceCount uint8) []uint16 {
	var result []uint16
	result = append(result, []uint16{0xb5, 0, 0x64, uint16(deviceStart), uint16(deviceCount)}...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}

//一键开合闸命令加密，参数stateArr的下标为设备号，值1表示开，0表示关
func allSwitchEncode(stateArr []int) []uint16 {
	var finalResult []uint16
	finalResult = append(finalResult, uint16(MSGTYPE_ALLSWITCH))
	data := make([]uint16, 64)
	length := len(stateArr)
	for i := 0; i < length/8+1; i++ {
		if length > 8 {
			data[i] = parseSwitchOperation(stateArr[i*8:i*8+8], 1)
			data[32+i] = parseSwitchOperation(stateArr[i*8:i*8+8], 0)
		} else {
			data[i] = parseSwitchOperation(stateArr[i*8:], 1)
			data[32+i] = parseSwitchOperation(stateArr[i*8:], 0)
		}
		length -= 8
	}
	finalResult = append(finalResult, data...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(finalResult)), CRC_BIGEND)
	finalResult = append(finalResult, byte2uint16(crc)...)
	return finalResult
}

//设备定时器设置——单次，响应解析使用writeRangeDecode，group只能是0-1-2-3
func deviceTimerOnce(addr, registerTable, group, timerType, state uint8, timestamp uint32) []uint16 {
	var registerNo uint16
	var result []uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_3
		}
	case REGISTER_TABLE_TWO:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_3
		}
	}
	result = append(result, []uint16{uint16(MSGTYPE_CMD), uint16(addr), 16}...)
	src := []uint16{registerNo, 8}
	dst := make([]byte, 16)
	for i := 0; i < 4; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	putUint16Switch(dst[4:6], uint16(group), BIGEND)
	putUint16Switch(dst[6:8], uint16(timerType), BIGEND)
	putUint16Switch(dst[8:10], 0, BIGEND)
	putUint16Switch(dst[10:12], uint16(state), BIGEND)
	putUint32Switch(dst[12:16], timestamp, BIGEND)
	result = append(result, byte2uint16(dst[0:4])...)
	result = append(result, 16)
	result = append(result, byte2uint16(dst[4:16])...)
	result = append(result, []uint16{0, 0, 0, 0}...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}

//设备定时器设置——每天，响应解析使用writeRangeDecode，group只能是0-1-2-3
func deviceTimerDaily(addr, registerTable, group, timerType, state, hour, minute uint8) []uint16 {
	var registerNo uint16
	var result []uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_3
		}
	case REGISTER_TABLE_TWO:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_3
		}
	}
	result = append(result, []uint16{uint16(MSGTYPE_CMD), uint16(addr), 16}...)
	src := []uint16{registerNo, 8}
	dst := make([]byte, 16)
	for i := 0; i < 4; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	putUint16Switch(dst[4:6], uint16(group), BIGEND)
	putUint16Switch(dst[6:8], uint16(timerType), BIGEND)
	putUint16Switch(dst[8:10], 1, BIGEND)
	putUint16Switch(dst[10:12], uint16(state), BIGEND)
	putUint16Switch(dst[12:14], uint16(hour), BIGEND)
	putUint16Switch(dst[14:16], uint16(minute), BIGEND)
	result = append(result, byte2uint16(dst[0:4])...)
	result = append(result, 16)
	result = append(result, byte2uint16(dst[4:16])...)
	result = append(result, []uint16{0, 0, 0, 0}...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}

//设备定时器设置——每周，响应解析使用writeRangeDecode，group只能是0-1-2-3
func deviceTimerWeekly(addr, registerTable, group, timerType, state, weekday, hour, minute uint8) []uint16 {
	var registerNo uint16
	var result []uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_3
		}
	case REGISTER_TABLE_TWO:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_3
		}
	}
	var realWeekday uint16 = 0
	//星期天是0
	if weekday == 7 {
		realWeekday = 0
	} else {
		realWeekday = 1 << weekday
	}
	result = append(result, []uint16{uint16(MSGTYPE_CMD), uint16(addr), 16}...)
	src := []uint16{registerNo, 8}
	dst := make([]byte, 18)
	for i := 0; i < 4; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	putUint16Switch(dst[4:6], uint16(group), BIGEND)
	putUint16Switch(dst[6:8], uint16(timerType), BIGEND)
	putUint16Switch(dst[8:10], 2, BIGEND)
	putUint16Switch(dst[10:12], uint16(state), BIGEND)
	putUint16Switch(dst[12:14], realWeekday, BIGEND)
	putUint16Switch(dst[14:16], uint16(hour), BIGEND)
	putUint16Switch(dst[16:18], uint16(minute), BIGEND)
	result = append(result, byte2uint16(dst[0:4])...)
	result = append(result, 16)
	result = append(result, byte2uint16(dst[4:18])...)
	result = append(result, []uint16{0, 0}...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}

//设备定时器设置——每月，响应解析使用writeRangeDecode，group只能是0-1-2-3
func deviceTimerMonthly(addr, registerTable, group, timerType, state, day, hour, minute uint8) []uint16 {
	var registerNo uint16
	var result []uint16
	switch registerTable {
	case REGISTER_TABLE_ONE:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_ONE_STTIMER_3
		}
	case REGISTER_TABLE_TWO:
		switch group {
		case 0:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_0
		case 1:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_1
		case 2:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_2
		case 3:
			registerNo = DEVICE_BASEADDR + DEVICEREGISTER_TWO_STTIMER_3
		}
	}
	result = append(result, []uint16{uint16(MSGTYPE_CMD), uint16(addr), 16}...)
	src := []uint16{registerNo, 8}
	dst := make([]byte, 20)
	for i := 0; i < 4; i += 2 {
		putUint16Switch(dst[i:i+2], src[i/2], BIGEND)
	}
	putUint16Switch(dst[4:6], uint16(group), BIGEND)
	putUint16Switch(dst[6:8], uint16(timerType), BIGEND)
	putUint16Switch(dst[8:10], 3, BIGEND)
	putUint16Switch(dst[10:12], uint16(state), BIGEND)
	putUint32Switch(dst[12:16], 1<<(day-1), BIGEND)
	putUint16Switch(dst[16:18], uint16(hour), BIGEND)
	putUint16Switch(dst[18:20], uint16(minute), BIGEND)
	result = append(result, byte2uint16(dst[0:4])...)
	result = append(result, 16)
	result = append(result, byte2uint16(dst[4:20])...)
	crc := make([]byte, 2)
	putUint16Switch(crc, crcGenerator(uint162byte(result)), CRC_BIGEND)
	result = append(result, byte2uint16(crc)...)
	return result
}
