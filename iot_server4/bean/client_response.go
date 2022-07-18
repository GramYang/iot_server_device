package bean

import (
	"reflect"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

const (
	DEVICE_ALL_REGISTER              = 200
	DEVICE_SNAPSHOOT                 = 201
	DEVICE_RESULT_MESSAGE            = 202
	GET_DEVICELIST_RESULT            = 203
	SWITCH_STATE_EVENT               = 204
	SWITCH_MODE_EVENT                = 205
	DEVICE_LOOP                      = 206
	SWITCH_LOOP                      = 207
	SWITCH_FAULT_EVENT               = 208
	SWITCH_WARN_EVENT                = 209
	SWITCH_HARDWARE_FAULT_EVENT      = 210
	SWITCH_ELECTRICLEAKAGE_TESTEVENT = 211
	MODBUS_SINGLE_WRITE_RESULT       = 212
	SWITCH_SETTING_UPLOAD            = 213
	DEVICE_ELECTRIC_QUANTITY_RESULT  = 214
	DEVICE_ALL_STATE_RESULT          = 215
	SHUTDOWN                         = 216
	RESULT_MESSAGE                   = 217
	CLIENT_HEARTBEAT_RESPONSE        = 218
	SWITCH_RUNTIME_RESPONSE          = 219
)

type DeviceAllRegister struct {
	DeviceId        string
	SignalIntensity int
	InternetMode    int
	DeviceCount     int
	UploadInterval  int
}

type DeviceSnapShoot struct {
	DeviceId    string
	DeviceCount int
	SnapShoots  []DeviceSnapShootContent1
}

type DeviceSnapShootContent1 struct {
	HardFault     int
	AlarmState    uint32
	ErrorState    uint32
	SwitchState   int
	SwitchMode    int
	Rms_u         float32
	Rms_i         float32
	Power_p       int
	Power_q       int
	Energy_pt     float32
	Energy_pt_l   float32
	Pf            float32
	Freq          float32
	Rms_IL        int
	Temperature   float32
	MotorRunTimes int
	ErrorTimes    int
}

type DeviceResultMessage struct {
	IsSuccess bool
	DeviceId  string
	Message   string
}

type GetDeviceListResult struct {
	DeviceIds []string
	Products  []string
}

type SwitchStateEvent struct {
	DeviceId     string
	Addr         int
	DeviceType   int
	RatedCurrent int
	SwitchStatus int
}

type SwitchModeEvent struct {
	DeviceId     string
	Addr         int
	DeviceType   int
	RatedCurrent int
	SwitchMode   int
}

type DeviceLoop struct {
	DeviceId    string
	DeviceCount int
	LeftTime    int
	List        []SwitchSnapShoot
}

type SwitchSnapShoot struct {
	Addr         int
	DeviceType   string
	RatedCurrent int
	PowerP       uint32
	Energy_pt    float32
}

type SwitchLoop struct {
	DeviceId     string
	Addr         int
	DeviceType   int
	RatedCurrent int
	LeftTime     int
	Runtime      SwitchRuntime1
}

type SwitchRuntime1 struct {
	HardFault     int
	AlarmState    uint32
	ErrorState    uint32
	SwitchState   int
	SwitchMode    int
	Rms_u         float32
	Rms_i         float32
	Power_p       int
	Power_q       int
	Energy_pt     float32
	Energy_pt_l   float32
	Pf            float32
	Freq          float32
	Rms_IL        int
	Temperature   float32
	MotorRunTimes int
	ErrorTimes    int
}

type SwitchFaultEvent struct {
	DeviceId     string
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	FaultEvents  uint32
	Data         FaultEventData
}

type FaultEventData struct {
	OverCurrent      float32 `byte:"4" division:"1000"`
	OverVoltage      float32 `byte:"2" division:"10"`
	UnderVoltage     float32 `byte:"2" division:"10"`
	Overdrive        uint32  `byte:"4"`
	ElectricQuantity float32 `byte:"4" division:"10"`
	OverHeat         float32 `byte:"2" division:"10"`
	ShortOut         float32 `byte:"4" division:"1000"`
	ElectricLeakage  uint16  `byte:"2"`
	GroundElectrode  uint16  `byte:"2"`
	Ack              uint16  `byte:"2"`
}

type SwitchWarnEvent struct {
	DeviceId     string
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	WarnEvents   uint32
	Data         WarnEventData
}

type WarnEventData struct {
	OverCurrent      float32 `byte:"4" division:"1000"`
	OverVoltage      float32 `byte:"2" division:"10"`
	UnderVoltage     float32 `byte:"2" division:"10"`
	Overdrive        uint32  `byte:"4"`
	ElectricQuantity float32 `byte:"4" division:"10"`
	OverHeat         float32 `byte:"2" division:"10"`
	ShortOut         float32 `byte:"4" division:"1000"`
	ElectricLeakage  uint16  `byte:"2"`
	GroundElectrode  uint16  `byte:"2"`
	Ack              uint16  `byte:"2"`
}

type SwitchHardwareFaultEvent struct {
	DeviceId       string
	Addr           uint8
	DeviceType     uint8
	RatedCurrent   uint16
	HardwareEvents uint16
}

type SwitchElectricLeakageTestEvent struct {
	DeviceId     string
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	Status       string
}

type ModbusSingleWriteResult struct {
	DeviceId      string
	Product       string
	Ok            bool
	Addr          int
	RegisterIndex int
	NewValue      int
	ErrorMsg      string
}

type SwitchSettingUpload struct {
	DeviceId     string
	Product      string
	Addr         int
	DeviceType   int
	RatedCurrent int
	Data         SwitchSetting1
	Timers       []SwitchTimer1
}

type SwitchSetting1 struct {
	AlarmEnable           uint16 `byte:"2"`
	ErrorEnable           uint16 `byte:"2"`
	VoltageLimitRstEnable uint16 `byte:"2"`
	IH_P                  uint16 `byte:"2"`
	IH                    uint16 `byte:"2"`
	UH_P                  uint16 `byte:"2"`
	UH                    uint16 `byte:"2"`
	UL_P                  uint16 `byte:"2"`
	UL                    uint16 `byte:"2"`
	PH_P                  uint16 `byte:"2"`
	PH                    uint16 `byte:"2"`
	EH_P                  uint32 `byte:"4"`
	EH                    uint32 `byte:"4"`
	IL_P                  uint16 `byte:"2"`
	IL                    uint16 `byte:"2"`
	TH_P                  uint16 `byte:"2"` //这个值必须在35-125内
	TH                    uint16 `byte:"2"` //这个值必须在35-125内
	UHL_CT                uint16 `byte:"2"` //这个值必须在0-60内
	UHL_RT                uint16 `byte:"2"` //这个值必须在10-60内
	IH_PH_CT              uint16 `byte:"2"`
}

type SwitchTimer1 struct {
	Group uint16
	Type  uint16 //任务
	Cycle uint16 //单次/天/周/月
	State uint16
	Timer interface{}
}

//一次性timer
type Once1 struct {
	Timestamp uint32
}

//每日timer
type Daily1 struct {
	Hour   uint16
	Minute uint16
}

//每周timer
type Weekly1 struct {
	WeekDay uint16 //1234567对应星期1到星期天
	Hour    uint16
	Minute  uint16
}

//每月timer
type Monthly1 struct {
	Day    uint16 //换算成了日期数
	Hour   uint16
	Minute uint16
}

//设备历史电量查询响应
type DeviceElectricQuantityResult struct {
	DeviceId   string
	Type       int       //1日2月3年
	PowerQ     []int     //W整数
	EnergyPt   []float32 //0.01度浮点
	RmsI       []float32 //0.001A浮点
	RecordTime int64
	Year       int
	Month      int
}

type DeviceAllStateResult struct {
	DeviceId    string
	DeviceCount int
	List        []SwitchState
}

type SwitchState struct {
	Addr         int
	RatedCurrent int
	IsOnLine     uint8
	HardFault    uint16
	AlarmState   uint32
	ErrorState   uint32
	SwitchState  uint16
	SwitchMode   uint16
}

//关闭前发送消息
type Shutdown struct {
	Code    int //1：token不同，其他用户登录。2：token过期。3：token错误。4：请求参数错误。5：服务器错误（创建user失败）
	Message string
}

//请求数据库结果消息，是同步的
type ResultMessage struct {
	IsSuccess bool
	Message   string
}

//心跳包响应
type ClientHeartBeatResponse struct {
	Username string
	Token    string
}

//单次查询的子开关runtime数据
type SwitchRuntimeResponse struct {
	DeviceId string
	Addr     int
	Data     SwitchRuntime1
}

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceAllRegister)(nil)).Elem(),
		ID:    DEVICE_ALL_REGISTER,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceSnapShoot)(nil)).Elem(),
		ID:    DEVICE_SNAPSHOOT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceResultMessage)(nil)).Elem(),
		ID:    DEVICE_RESULT_MESSAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*GetDeviceListResult)(nil)).Elem(),
		ID:    GET_DEVICELIST_RESULT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchStateEvent)(nil)).Elem(),
		ID:    SWITCH_STATE_EVENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchModeEvent)(nil)).Elem(),
		ID:    SWITCH_MODE_EVENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceLoop)(nil)).Elem(),
		ID:    DEVICE_LOOP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchLoop)(nil)).Elem(),
		ID:    SWITCH_LOOP,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchFaultEvent)(nil)).Elem(),
		ID:    SWITCH_FAULT_EVENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchWarnEvent)(nil)).Elem(),
		ID:    SWITCH_WARN_EVENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchHardwareFaultEvent)(nil)).Elem(),
		ID:    SWITCH_HARDWARE_FAULT_EVENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchElectricLeakageTestEvent)(nil)).Elem(),
		ID:    SWITCH_ELECTRICLEAKAGE_TESTEVENT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*ModbusSingleWriteResult)(nil)).Elem(),
		ID:    MODBUS_SINGLE_WRITE_RESULT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchSettingUpload)(nil)).Elem(),
		ID:    SWITCH_SETTING_UPLOAD,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceElectricQuantityResult)(nil)).Elem(),
		ID:    DEVICE_ELECTRIC_QUANTITY_RESULT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*DeviceAllStateResult)(nil)).Elem(),
		ID:    DEVICE_ALL_STATE_RESULT,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*Shutdown)(nil)).Elem(),
		ID:    SHUTDOWN,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*ResultMessage)(nil)).Elem(),
		ID:    RESULT_MESSAGE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*ClientHeartBeatResponse)(nil)).Elem(),
		ID:    CLIENT_HEARTBEAT_RESPONSE,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("json"),
		Type:  reflect.TypeOf((*SwitchRuntimeResponse)(nil)).Elem(),
		ID:    SWITCH_RUNTIME_RESPONSE,
	})
}
