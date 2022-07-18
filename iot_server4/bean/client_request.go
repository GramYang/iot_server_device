package bean

import (
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"

	_ "github.com/davyxu/cellnet/codec/json"
)

const (
	CLIENT_HEARTBEAT             = 3
	GET_DEVICELIST               = 4
	MULTI_SWITCH_OPERATION       = 5
	SWITCH_ON                    = 6
	SWITCH_OFF                   = 7
	SWITCH_LOCK_ON               = 8
	SWITCH_LOCK_OFF              = 9
	DEVICE_RESET                 = 10
	SET_TIMER                    = 11
	SWITCH_ELECTRIC_LEAKAGE_TEST = 12
	SWITCH_ALARM_ENABLE          = 13
	SWITCH_ERROR_ENABLE          = 14
	DEVICE_ELECTRIC_QUANTITY     = 15
	DEVICE_ALL_STATE             = 16
	GET_SWITCH_SETTING           = 17
	SWITCH_LOOP_ON               = 18
	SWITCH_LOOP_OFF              = 19
	SWITCH_CLEAR_FAULT           = 20
	VOLTAGE_LIMTRST_ENABLE       = 21
	IH_P                         = 22
	IH                           = 23
	UH_P                         = 24
	UH                           = 25
	UL_P                         = 26
	UL                           = 27
	PH_P                         = 28
	PH                           = 29
	EH_P                         = 30
	EH                           = 31
	IL_P                         = 32
	IL                           = 33
	TH_P                         = 34
	TH                           = 35
	UHL_CT                       = 36
	UHL_RT                       = 37
	IH_PH_CT                     = 38
	GET_ALL_SWITCH_RUNTIME       = 39
	CLIENT_LOGOUT                = 40
)

//测试消息
type TestEchoACK struct {
	Msg   string
	Value int32
}

//客户端心跳
type ClientHeartBeat struct {
	Username string
	Token    string
}

//请求用户绑定设备列表
type GetDeviceList struct {
	Heartbeat ClientHeartBeat
}

//多开关一键开合
type MultiSwitchOperation struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Cmd       []int
}

//开关开闸
type SwitchOn struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//开关关闸
type SwitchOff struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//开关锁定
type SwitchLockOn struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//开关解锁
type SwitchLockOff struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//设备重分配地址
type DeviceReset struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
}

//设置定时器
type SetTimer struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Group     int
	Task      int
	Cycle     int
	State     int
	Timestamp uint32
	Minute    int
	Hour      int
	WeekDay   int //1-7是星期1到星期天
	Day       int //1-31号
}

//开关漏电测试
type SwitchElectricLeakageTest struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//故障预警开启
type SwitchAlarmEnable struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Enables   []int
}

//故障保护开启
type SwitchErrorEnable struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Enables   []int
}

//设置历史电量查询
type DeviceElectricQuantity struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Type      int //1日2月3年
	Day       int
	Month     int
	Year      int
}

type DeviceAllState struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
}

type GetSwitchSetting struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//开启子开关轮询
type SwitchLoopOn struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//关闭子开关轮询
type SwitchLoopOff struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//清除开关故障
type SwitchClearFault struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
}

//过欠压恢复使能
type VoltageLimitRstEnable struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Enable    bool
}

type SetIHP struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}

type SetIH struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetUHP struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetUH struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetULP struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetUL struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetPHP struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetPH struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetEHP struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     uint32
}
type SetEH struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Index     int //不能是0
	Value     int //这里有可能是负值
}
type SetILP struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetIL struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetTHP struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetTH struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetUHLCT struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}

type SetUHLRT struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}
type SetIHPHCT struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
	Addr      int
	Value     int
}

type GetAllSwitchRuntime struct {
	Heartbeat ClientHeartBeat
	DeviceId  string
	Product   string
}

type ClientLogout struct {
	Heartbeat ClientHeartBeat
}

func init() {
	//测试请求bean
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    999,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*ClientHeartBeat)(nil)).Elem(),
		ID:    CLIENT_HEARTBEAT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*GetDeviceList)(nil)).Elem(),
		ID:    GET_DEVICELIST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*MultiSwitchOperation)(nil)).Elem(),
		ID:    MULTI_SWITCH_OPERATION,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchOn)(nil)).Elem(),
		ID:    SWITCH_ON,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchOff)(nil)).Elem(),
		ID:    SWITCH_OFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchLockOn)(nil)).Elem(),
		ID:    SWITCH_LOCK_ON,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchLockOff)(nil)).Elem(),
		ID:    SWITCH_LOCK_OFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceReset)(nil)).Elem(),
		ID:    DEVICE_RESET,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetTimer)(nil)).Elem(),
		ID:    SET_TIMER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchElectricLeakageTest)(nil)).Elem(),
		ID:    SWITCH_ELECTRIC_LEAKAGE_TEST,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchAlarmEnable)(nil)).Elem(),
		ID:    SWITCH_ALARM_ENABLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchErrorEnable)(nil)).Elem(),
		ID:    SWITCH_ERROR_ENABLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceElectricQuantity)(nil)).Elem(),
		ID:    DEVICE_ELECTRIC_QUANTITY,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceAllState)(nil)).Elem(),
		ID:    DEVICE_ALL_STATE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*GetSwitchSetting)(nil)).Elem(),
		ID:    GET_SWITCH_SETTING,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchLoopOn)(nil)).Elem(),
		ID:    SWITCH_LOOP_ON,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchLoopOff)(nil)).Elem(),
		ID:    SWITCH_LOOP_OFF,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchClearFault)(nil)).Elem(),
		ID:    SWITCH_CLEAR_FAULT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*VoltageLimitRstEnable)(nil)).Elem(),
		ID:    VOLTAGE_LIMTRST_ENABLE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetIHP)(nil)).Elem(),
		ID:    IH_P,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetIH)(nil)).Elem(),
		ID:    IH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetUHP)(nil)).Elem(),
		ID:    UH_P,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetUH)(nil)).Elem(),
		ID:    UH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetULP)(nil)).Elem(),
		ID:    UL_P,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetUL)(nil)).Elem(),
		ID:    UL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetPHP)(nil)).Elem(),
		ID:    PH_P,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetPH)(nil)).Elem(),
		ID:    PH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetEHP)(nil)).Elem(),
		ID:    EH_P,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetEH)(nil)).Elem(),
		ID:    EH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetILP)(nil)).Elem(),
		ID:    IL_P,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetIL)(nil)).Elem(),
		ID:    IL,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetTHP)(nil)).Elem(),
		ID:    TH_P,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetTH)(nil)).Elem(),
		ID:    TH,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetUHLCT)(nil)).Elem(),
		ID:    UHL_CT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetUHLRT)(nil)).Elem(),
		ID:    UHL_RT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SetIHPHCT)(nil)).Elem(),
		ID:    IH_PH_CT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*GetAllSwitchRuntime)(nil)).Elem(),
		ID:    GET_ALL_SWITCH_RUNTIME,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*ClientLogout)(nil)).Elem(),
		ID:    CLIENT_LOGOUT,
	})
}
