package modbus

import "errors"

//modbus——设备&主机寄存器单独写响应解析
func WriteSingleDecode(data []uint16) *WriteSingleResult {
	result := &WriteSingleResult{}
	result.DeviceAddr = uint8(data[0])
	if data[1] == WRITESINGLE_SUCCESS {
		result.Ok = true
		result.RegisterIndex = uint16Switch(uint162byte(data[2:4]), BIGEND)
		result.NewValue = uint16Switch(uint162byte(data[4:6]), BIGEND)
	} else if data[1] == WRITESINGLE_FAILED {
		result.Ok = false
		result.ErrorMsg = ERRORCODE_MSG[uint8(data[2])]
	}
	return result
}

//modbus——设备&主机寄存器范围写响应解析
func WriteRangeDecode(data []uint16) *WriteRangeResult {
	result := &WriteRangeResult{}
	result.DeviceAddr = uint8(data[0])
	if data[1] == WRITERANGE_SUCCESS {
		result.Ok = true
		result.RegisterStartIndex = uint16Switch(uint162byte(data[2:4]), BIGEND)
		result.RegisterCount = uint16Switch(uint162byte(data[4:6]), BIGEND)
	} else if data[1] == WRITERANGE_FAILED {
		result.Ok = false
		result.ErrorMsg = ERRORCODE_MSG[uint8(data[2])]
	}
	return result
}

//modbus——解析设备&主机寄存器范围读响应，ptr是数据bean的实例指针
func ReadRangeDecode(data []uint16) *ReadRangeResult {
	result := &ReadRangeResult{}
	result.DeviceAddr = uint8(data[0])
	if data[1] == READRANGE_SUCCESS {
		result.Ok = true
		result.DataSize = uint8(data[2])
	} else if data[1] == READRANGE_FAILED {
		result.Ok = false
		result.ErrorMsg = ERRORCODE_MSG[uint8(data[2])]
	}
	return result
}

func ReadRangeSuccessDecodeWithPtr(result *ReadRangeResult, data []uint16, ptr interface{}) {
	intSliceMapBean(data[3:3+data[2]], ptr)
	if result != nil {
		result.Data = ptr
	}
}

//modbus——解析主机快照范围读响应——成功
func snapShootReadRangeSuccessDecode(data []uint16) *SnapShootReadRangeResult {
	result := &SnapShootReadRangeResult{}
	if data[1] == SNAPSHOOT_READRANGE_SUCCESS {
		result.Ok = true
		result.DeviceCount = uint8(data[2])
		var index = 3
		for i := 0; i < int(data[2]); i++ {
			snapShoot := DeviceSnapShoot{}
			snapShoot.DeviceAddr = uint8(data[index])
			snapShoot.DeviceType = uint8(data[index+1])
			if checkDeviceType(snapShoot.DeviceType) == REGISTER_TABLE_ONE {
				snapShoot1 := DeviceSnapShootContent1{}
				intSliceMapBean(data[index+2:index+36], &snapShoot1)
				snapShoot.Data = snapShoot1
			}
			result.DeviceSnapShoots = append(result.DeviceSnapShoots, snapShoot)
			index += 36
		}
	} else if data[1] == SNAPSHOOT_READRANGE_FAILED {
		result.Ok = false
		result.ErrorMsg = ERRORCODE_MSG[uint8(data[2])]
	}
	return result
}

//主机快照定时上传解析，和上面比起来少了主机地址和功能码
func SnapShootTimedUploadDecode(data []uint16) *SnapShootReadRangeResult {
	result := &SnapShootReadRangeResult{}
	result.Ok = true
	result.DeviceCount = uint8(data[0])
	var index = 1
	for i := 0; i < int(data[0]); i++ {
		snapShoot := DeviceSnapShoot{}
		snapShoot.DeviceAddr = uint8(data[index])
		snapShoot.DeviceType = uint8(data[index+1])
		snapShoot.RatedCurrent = uint16Switch(uint162byte(data[index+2:index+4]), BIGEND)
		if checkDeviceType(snapShoot.DeviceType) == REGISTER_TABLE_ONE {
			snapShoot1 := DeviceSnapShootContent1{}
			intSliceMapBean(data[index+4:index+52], &snapShoot1)
			snapShoot.Data = snapShoot1
			result.DeviceSnapShoots = append(result.DeviceSnapShoots, snapShoot)
			index += 52
		} else if checkDeviceType(snapShoot.DeviceType) == REGISTER_TABLE_TWO {
			snapShoot2 := DeviceSnapShootContent2{}
			intSliceMapBean(data[index+4:index+56], &snapShoot2)
			snapShoot.Data = snapShoot2
			result.DeviceSnapShoots = append(result.DeviceSnapShoots, snapShoot)
			index += 56
		}
	}
	return result
}

func SnapShootMultiUploadDecode(data []uint16) *SnapShootResult {
	result := &SnapShootResult{}
	result.LeftTime = uint8(data[0])
	result.DeviceCount = uint8(data[1])
	var index = 2
	for i := 0; i < int(data[1]); i++ {
		snapShoot := DeviceSnapShoot{}
		snapShoot.DeviceAddr = uint8(data[index])
		snapShoot.DeviceType = uint8(data[index+1])
		snapShoot.RatedCurrent = uint16Switch(uint162byte(data[index+2:index+4]), BIGEND)
		snapShootContent3 := DeviceSnapShootContent3{}
		intSliceMapBean(data[index+4:index+12], &snapShootContent3)
		snapShoot.Data = snapShootContent3
		result.DeviceSnapShoots = append(result.DeviceSnapShoots, snapShoot)
		index += 12
	}
	return result
}

//一键开合闸响应解析
func allSwitchDecode(data []uint16) ([]int, error) {
	if len(data) != 32 {
		return nil, errors.New("data length error")
	}
	var res []int
	for _, v := range data {
		arr := parseSwitchState(v)
		res = append(res, arr...)
	}
	return res, nil
}

//单个定时器状态响应解析
func deviceTimerDecode1(data []uint16) *DeviceTimer1 {
	tmp := &DeviceTimer1{}
	tmp.Group = uint16Switch(uint162byte(data[0:2]), BIGEND)
	tmp.Type = uint16Switch(uint162byte(data[2:4]), BIGEND)
	tmp.Cycle = uint16Switch(uint162byte(data[4:6]), BIGEND)
	tmp.State = uint16Switch(uint162byte(data[6:8]), BIGEND)
	switch tmp.Cycle {
	case 0:
		tmp.Timer = Once1{Timestamp: uint32Switch(uint162byte(data[8:12]), BIGEND)}
	case 1:
		tmp.Timer = Daily1{Hour: uint16Switch(uint162byte(data[8:10]), BIGEND), Minute: uint16Switch(uint162byte(data[10:12]), BIGEND)}
	case 2:
		tmp.Timer = Weekly1{WeekDay: parseWeekDay(uint16Switch(uint162byte(data[8:10]), BIGEND)), Hour: uint16Switch(uint162byte(data[10:12]), BIGEND), Minute: uint16Switch(uint162byte(data[12:14]), BIGEND)}
	case 3:
		tmp.Timer = Monthly1{Day: parseMonthDay(data[8:12]), Hour: uint16Switch(uint162byte(data[12:14]), BIGEND), Minute: uint16Switch(uint162byte(data[14:16]), BIGEND)}
	}
	return tmp
}
