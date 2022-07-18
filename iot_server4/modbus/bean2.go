package modbus

type DeviceSnapShootContent2 struct {
	StartReturn   uint32  `byte:"4"`
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

type DeviceRegisterDataRead2 struct {
	StartReturn   uint32  `byte:"4"`
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
