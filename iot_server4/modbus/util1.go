package modbus

import (
	"encoding/binary"
	"reflect"
	"strconv"
	"strings"
)

//************大小端切换封装*****************
func putUint16Switch(bs []byte, v uint16, isBigEnd bool) {
	if isBigEnd {
		binary.BigEndian.PutUint16(bs, v)
	} else {
		binary.LittleEndian.PutUint16(bs, v)
	}
}

func uint16Switch(bs []byte, isBigEnd bool) uint16 {
	if isBigEnd {
		return binary.BigEndian.Uint16(bs)
	} else {
		return binary.LittleEndian.Uint16(bs)
	}
}

func putUint32Switch(bs []byte, v uint32, isBigEnd bool) {
	if isBigEnd {
		binary.BigEndian.PutUint32(bs, v)
	} else {
		binary.LittleEndian.PutUint32(bs, v)
	}
}

func uint32Switch(bs []byte, isBigEnd bool) uint32 {
	if isBigEnd {
		return binary.BigEndian.Uint32(bs)
	} else {
		return binary.LittleEndian.Uint32(bs)
	}
}

//************大小端切换封装*****************

//生成crc16
func crcGenerator(bs []byte) uint16 {
	var uchCRCHi uint8 = 0xff
	var uchCRCLo uint8 = 0xff
	var uIndex uint8
	for _, v := range bs {
		uIndex = uchCRCHi ^ (v)
		uchCRCHi = uchCRCLo ^ auchCRCHi[uIndex]
		uchCRCLo = auchCRCLo[uIndex]
	}
	return uint16(uchCRCLo)<<8 | uint16(uchCRCHi)
}

//验证crc16
func CheckCRC(data []uint16) bool {
	return uint16Switch(uint162byte(data[len(data)-2:]), CRC_BIGEND) == crcGenerator(uint162byte(data[:len(data)-2]))
}

//byte转uint16
func byte2uint16(bs []byte) []uint16 {
	if len(bs) == 0 {
		return nil
	}
	var us []uint16
	for _, v := range bs {
		us = append(us, uint16(v))
	}
	return us
}

//uint16转byte
func uint162byte(us []uint16) []byte {
	if len(us) == 0 {
		return nil
	}
	var bs []byte
	for _, v := range us {
		bs = append(bs, byte(v))
	}
	return bs
}

//int转uint16
func Int2uint16(is []int) []uint16 {
	if len(is) == 0 {
		return nil
	}
	var us []uint16
	for _, v := range is {
		us = append(us, uint16(v))
	}
	return us
}

//data是待解析的数据，ptr是一个数据bean实例指针
func intSliceMapBean(data []uint16, ptr interface{}) {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	var totalSize int
	var propSize []int
	for i := 0; i < typ.NumField(); i++ {
		tmp, _ := strconv.Atoi(typ.Field(i).Tag.Get("byte"))
		propSize = append(propSize, tmp)
		totalSize += tmp
	}
	if totalSize != len(data) {
		return
	}
	var index int
	for i := 0; i < typ.NumField(); i++ {
		switch typ.Field(i).Type.Kind() {
		case reflect.Uint8:
			if propSize[i] == 1 {
				val.Field(i).Set(reflect.ValueOf(uint8(data[index])))
				index += 1
			}
		case reflect.Uint16:
			if propSize[i] == 2 {
				bs := uint162byte(data[index : index+2])
				val.Field(i).Set(reflect.ValueOf(uint16Switch(bs, BIGEND)))
				index += 2
			}
		case reflect.Int16:
			if propSize[i] == 2 {
				bs := uint162byte(data[index : index+2])
				tmp := int16(uint16Switch(bs, BIGEND))
				val.Field(i).Set(reflect.ValueOf(tmp))
				index += 2
			}
		case reflect.Float32:
			if propSize[i] == 2 {
				bs := uint162byte(data[index : index+2])
				tmp := int16(uint16Switch(bs, BIGEND))
				div, _ := strconv.Atoi(typ.Field(i).Tag.Get("division"))
				val.Field(i).Set(reflect.ValueOf(float32(tmp) / float32(div))) //这里除法没有误差，但是乘法就有误差了，需要用decimal库来操作
				index += 2
			} else if propSize[i] == 4 {
				bs := uint162byte(data[index : index+4]) //4字长值传输时非标准大端
				tmp := uint32Switch(bs, BIGEND)
				div, _ := strconv.Atoi(typ.Field(i).Tag.Get("division"))
				val.Field(i).Set(reflect.ValueOf(float32(tmp) / float32(div)))
				index += 4
			}
		case reflect.Uint32:
			if propSize[i] == 4 {
				bs := uint162byte(data[index : index+4]) //4字长值传输时非标准大端
				tmp := uint32Switch(bs, BIGEND)
				val.Field(i).Set(reflect.ValueOf(tmp))
				index += 4
			}
		}
	}
}

//设备定时器解析星期天数
func parseWeekDay(b uint16) uint16 {
	if b == 0 {
		return 7
	}
	for i := 0; i < 6; i++ {
		if b == 1<<i {
			return uint16(i + 1)
		}
	}
	return 0
}

//设备定时器解析月天数
func parseMonthDay(bs []uint16) uint16 {
	if len(bs) != 4 {
		return 0
	}
	v := uint32Switch(uint162byte(bs), BIGEND)
	for i := 0; i < 31; i++ {
		if v == 1<<i {
			return uint16(i + 1)
		}
	}
	return 0
}

//因为go环境的大小端解析只能使用无符号数，因此用无符号数来表示负数
func uint162int16(v uint16) uint16 {
	return (^(v << 1) >> 1) | (1 << 15) + 1
}

//bit偏移转换为值
func bitToData(bit uint8) uint16 {
	return 1 << bit
}

//多bit位转换为值
func bitsToData(bits []uint8) uint16 {
	var data uint16
	for _, v := range bits {
		data += 1 << v
	}
	return data
}

//解析uint32多位信息（uint16也可以用）
func parseMultiBitsOfUint32(bits uint32, msgMap map[uint8]string) []string {
	var result []string
	for k, v := range msgMap {
		if bits>>k&1 == 1 {
			result = append(result, v)
		}
	}
	return result
}

//格式化输出数组
func sliceUint16Fmt(data []uint16) string {
	var b strings.Builder
	b.WriteString("[")
	for k, v := range data {
		b.WriteString(strconv.Itoa(int(v)))
		if k != len(data)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("]")
	return b.String()
}

//返回设备类型对应的寄存器表号
func checkDeviceType(t uint8) uint8 {
	switch t {
	case 0, 1, 2, 3, 4, 5, 6, 7:
		return 1
	case 8, 9, 10, 11, 12, 13, 14, 15:
		return 2
	}
	return 0
}

//一键分合闸，每个长度为8的数组转换成一个值，flag为1表示开，0表示关
func parseSwitchOperation(arr []int, flag int) uint16 {
	var res uint16
	for k, v := range arr {
		if v == flag {
			res += 1 << k
		}
	}
	return res
}

//一键分合闸，根据一个值获取其代表8个开关的开启关闭状态
func parseSwitchState(v uint16) []int {
	res := make([]int, 8)
	for i := 0; i < 8; i++ {
		if v>>i&1 == 1 {
			res[i] = 1
		}
	}
	return res
}
