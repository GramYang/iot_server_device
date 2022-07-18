package modbus

//主机下设备快照响应
type SnapShootReadRangeResult struct {
	Ok               bool
	DeviceCount      uint8
	DeviceSnapShoots []DeviceSnapShoot
	ErrorMsg         string
}

type SnapShootResult struct {
	LeftTime         uint8
	DeviceCount      uint8
	DeviceSnapShoots []DeviceSnapShoot
}

type DeviceSnapShoot struct {
	DeviceAddr   uint8
	DeviceType   uint8
	RatedCurrent uint16
	Data         interface{}
}

type DeviceSnapShootContent1 struct {
	HardFault     uint16  `byte:"2"`
	AlarmState    uint32  `byte:"4"`
	ErrorState    uint32  `byte:"4"`
	SwitchState   uint16  `byte:"2"`
	SwitchMode    uint16  `byte:"2"`
	Rms_u         float32 `byte:"2" division:"10"`
	Rms_i         float32 `byte:"4" division:"1000"`
	Power_p       uint32  `byte:"4"`
	Power_q       uint32  `byte:"4"`
	Energy_pt     float32 `byte:"4" division:"100"`
	Energy_pt_l   float32 `byte:"4" division:"100"`
	Pf            float32 `byte:"2" division:"100"`
	Freq          float32 `byte:"2" division:"10"`
	Rms_IL        uint16  `byte:"2"`
	Temperature   float32 `byte:"2" division:"10"`
	MotorRunTimes uint16  `byte:"2"`
	ErrorTimes    uint16  `byte:"2"`
}

//设备&主机寄存器范围读响应
type ReadRangeResult struct {
	Ok         bool
	DeviceAddr uint8
	DataSize   uint8
	Data       interface{} //这里是指针类型
	ErrorMsg   string
}

//主机全寄存器
type HostRegisterData struct {
	HostVersion     uint16 `byte:"2"`
	SignalIntensity int16  `byte:"2"`
	InternetMode    uint16 `byte:"2"`
	DeviceCount     uint16 `byte:"2"`
	UploadInterval  uint16 `byte:"2"`
}

//设备寄存器——读取字段数据
type DeviceRegisterDataRead1 struct {
	HardFault     uint16  `byte:"2"`
	AlarmState    uint32  `byte:"4"`
	ErrorState    uint32  `byte:"4"`
	SwitchState   uint16  `byte:"2"`
	SwitchMode    uint16  `byte:"2"`
	Rms_U         float32 `byte:"2" division:"10"`
	Rms_I         float32 `byte:"4" division:"1000"`
	PowerP        uint32  `byte:"4"`
	PowerQ        uint32  `byte:"4"`
	Energy_pt     float32 `byte:"4" division:"100"`
	Energy_pt_l   float32 `byte:"4" division:"100"`
	Pf            float32 `byte:"2" division:"100"`
	Freq          float32 `byte:"2" division:"10"`
	Rms_il        uint16  `byte:"2"`
	Temperature   float32 `byte:"2" division:"10"`
	MotorRunTimes uint16  `byte:"2"`
	ErrorTimes    uint16  `byte:"2"`
}

//设备寄存器——设置字段数据
type DeviceRegisterDataWrite1 struct {
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

//设备寄存器——定时器数据
type DeviceTimer1 struct {
	Group uint16
	Type  uint16
	Cycle uint16
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
	WeekDay uint16
	Hour    uint16
	Minute  uint16
}

//每月timer
type Monthly1 struct {
	Day    uint16
	Hour   uint16
	Minute uint16
}

//设备&主机寄存器范围写响应
type WriteRangeResult struct {
	Ok                 bool
	DeviceAddr         uint8
	RegisterStartIndex uint16
	RegisterCount      uint16
	ErrorMsg           string
}

//modbus单独写响应
type WriteSingleResult struct {
	Ok            bool
	DeviceAddr    uint8
	RegisterIndex uint16
	NewValue      uint16
	ErrorMsg      string
}

//设备实时数据上报
type DeviceStatusUpload struct {
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	LeftTime     uint8
	Data         interface{}
}

//设备上报设置数据
type DeviceSettingStatus struct {
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	Data         interface{}
	Timers       interface{} //设备定时器默认全部启用，你只能开启和关闭
}

//设备故障事件上报
type DeviceFaultEvent struct {
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	FaultEvents  uint32
	Data         interface{}
}

type FaultData1 struct {
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

//设备预警事件上报
type DeviceWarnEvent struct {
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	WarnEvents   uint32
	Data         interface{}
}

type WarnData1 struct {
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

//设备硬件故障事件上报
type DeviceHardwareFaultEvent struct {
	Addr           uint8
	DeviceType     uint8
	RatedCurrent   uint16
	HardwareEvents uint16
}

//开关状态事件上报
type SwitchStatusEvent struct {
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	SwitchStatus uint16
}

//开关模式事件上报
type SwitchModeEvent struct {
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	SwitchMode   uint16
}

//漏电检测事件上报
type ElectricLeakageTestEvent struct {
	Addr         uint8
	DeviceType   uint8
	RatedCurrent uint16
	Status       string
}

//主机下设备状态: 读取字段的前5个
type DeviceState struct {
	IsOnLine    uint8  `byte:"1"`
	HardFault   uint16 `byte:"2"`
	AlarmState  uint32 `byte:"4"`
	ErrorState  uint32 `byte:"4"`
	SwitchState uint16 `byte:"2"`
	SwitchMode  uint16 `byte:"2"`
}

//所有设备类型一致，只含有功功率和总有功电量（8个字节）
type DeviceSnapShootContent3 struct {
	PowerP    uint32  `byte:"4"`
	Energy_pt float32 `byte:"4" division:"100"`
}
