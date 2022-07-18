package router

import (
	"iot_server4/bean"
	"iot_server4/command"
	"iot_server4/config"
	"iot_server4/modbus"
	"iot_server4/model"
	"time"

	g "github.com/GramYang/gylog"
	"github.com/davyxu/cellnet"
)

//寄存器表1的路由
func Router1(deviceId, product, timestamp string, data []int) {
	if len(data) <= 2 {
		g.Errorln("aliyun message empty")
		return
	}
	data16 := modbus.Int2uint16(data)
	if !modbus.CheckCRC(data16) {
		g.Errorln("aliyun message crc verification error")
		return
	}

	//更新deviceSession用于设备心跳，上线则执行回调
	model.UpdateDeviceSession(deviceId, time.Now().UnixMilli(), func(deviceId string) {
		g.Debugln("执行UpdateDeviceSession设备上线回调")
		model.SingleSend(deviceId, func(product string) {
			g.Debugf("设备 %s 全寄存器读\n", deviceId)
			command.HostAllRegister(deviceId, product) //主机全寄存器读，主机全寄存器读的响应会触发子开关轮询
		}, config.Conf.CmdInterval)
	})

	switch data16[0] {

	default:
		g.Debugln("未知响应")

	//响应时间戳请求，校正定时器
	case uint16(modbus.MSGTYPE_TIMESTAMP):
		g.Debugln("主机时间戳请求:")
		// model.SingleSend(deviceId, func(product string) {
		// 	command.HostTimestampResponse(deviceId, product)
		// }, config.Conf.CmdInterval)
		command.HostTimestampResponse(deviceId, product)

	//iccid
	case uint16(modbus.MSGTYPE_ICCID_RESULT):
		g.Debugln("主机iccid上报:")
		// result := modbus.HostIccidResult(data16[1 : len(data16)-2])

	//主机全寄存器上报
	case uint16(modbus.MSGTYPE_HOSTDETAIL):
		g.Debugln("主机全寄存器上报:")
		result := modbus.HostInfoEventUploadResult(data16[13:])
		var deviceNum = int(result.DeviceCount)
		//此消息在设备重置后会上电时主动上传，用来触发子开关轮询变更
		//管理开关实时数据轮询
		if deviceNum != 0 {
			model.UpdateDeviceCacheSwitchNumAndSend(deviceId, deviceNum, func(i int) {
				g.Debugf("关闭设备 %s 开关 %d 轮询\n", deviceId, i)
				command.SwitchRuntime(deviceId, product, uint8(i), uint8(config.Conf.LoopInterval), 0)
			}, config.Conf.CmdInterval)
			//loop更新设备子开关数，第一次赋值是就立即触发一次全子开关runtime读
			model.UpdateLoopCacheSwitchNumAndCall(deviceId, deviceNum, nil)
			var t1 []time.Duration
			var f1 []func()
			for i := 1; i <= deviceNum; i++ {
				t := i
				t1 = append(t1, config.Conf.CmdInterval)
				f1 = append(f1, func() {
					g.Debugf("单次请求设备 %s 子开关 %d runtime\n", deviceId, t)
					command.SwitchRuntimeOnce(deviceId, product, uint8(t))
				})
			}
			model.MultiSend(deviceId, f1, t1)
			//返回给客户端
			model.Broadcast(deviceId, func(ses cellnet.Session) {
				if ses == nil {
					return
				}
				ses.Send(&bean.DeviceAllRegister{
					DeviceId:        deviceId,
					SignalIntensity: int(result.SignalIntensity),
					InternetMode:    int(result.SignalIntensity),
					DeviceCount:     int(result.DeviceCount),
					UploadInterval:  int(result.UploadInterval),
				})
			})
		} else {
			//按道理说主机不会返回子开关数量为空，先留着吧。
		}

	//主机下全设备上报
	case uint16(modbus.MSGTYPE_ALLSTATE_RESULT):
		g.Debugln("主机下全设备上报:")
		result := modbus.ReadHostAllStateResult(data16[1 : len(data16)-2])
		var list []bean.SwitchState
		for _, v := range result.DeviceSnapShoots {
			tmp := bean.SwitchState{
				Addr:         int(v.DeviceAddr),
				RatedCurrent: int(v.RatedCurrent),
				IsOnLine:     v.Data.(modbus.DeviceState).IsOnLine,
				HardFault:    v.Data.(modbus.DeviceState).HardFault,
				AlarmState:   v.Data.(modbus.DeviceState).AlarmState,
				ErrorState:   v.Data.(modbus.DeviceState).ErrorState,
				SwitchState:  v.Data.(modbus.DeviceState).SwitchState,
				SwitchMode:   v.Data.(modbus.DeviceState).SwitchMode,
			}
			list = append(list, tmp)
		}
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.DeviceAllStateResult{
				DeviceId:    deviceId,
				DeviceCount: int(result.DeviceCount),
				List:        list,
			})
		})

	//主机快照定时上传
	case uint16(modbus.MSGTYPE_SNAPSHOOT_TIMEDUPLOAD):
		g.Debugln("主机快照定时:")
		result := modbus.SnapShootTimedUploadDecode(data16[1:])
		//发回给客户端
		dss := bean.DeviceSnapShoot{DeviceId: deviceId, DeviceCount: int(result.DeviceCount)}
		var dssc []bean.DeviceSnapShootContent1
		for _, v := range result.DeviceSnapShoots {
			dssc = append(dssc, bean.DeviceSnapShootContent1{
				HardFault:     int(v.Data.(modbus.DeviceSnapShootContent1).HardFault),
				AlarmState:    v.Data.(modbus.DeviceSnapShootContent1).AlarmState,
				ErrorState:    v.Data.(modbus.DeviceSnapShootContent1).ErrorState,
				SwitchState:   int(v.Data.(modbus.DeviceSnapShootContent1).SwitchState),
				SwitchMode:    int(v.Data.(modbus.DeviceSnapShootContent1).SwitchMode),
				Rms_u:         v.Data.(modbus.DeviceSnapShootContent1).Rms_u,
				Rms_i:         v.Data.(modbus.DeviceSnapShootContent1).Rms_i,
				Power_p:       int(v.Data.(modbus.DeviceSnapShootContent1).Power_p),
				Power_q:       int(v.Data.(modbus.DeviceSnapShootContent1).Power_q),
				Energy_pt:     v.Data.(modbus.DeviceSnapShootContent1).Energy_pt,
				Energy_pt_l:   v.Data.(modbus.DeviceSnapShootContent1).Energy_pt_l,
				Pf:            v.Data.(modbus.DeviceSnapShootContent1).Pf,
				Freq:          v.Data.(modbus.DeviceSnapShootContent1).Freq,
				Rms_IL:        int(v.Data.(modbus.DeviceSnapShootContent1).Rms_IL),
				Temperature:   v.Data.(modbus.DeviceSnapShootContent1).Temperature,
				MotorRunTimes: int(v.Data.(modbus.DeviceSnapShootContent1).MotorRunTimes),
				ErrorTimes:    int(v.Data.(modbus.DeviceSnapShootContent1).ErrorTimes),
			})
		}
		dss.SnapShoots = dssc
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(dss)
		})

	//主机轮询，由于武汉那边的傻逼不会用导致主机轮询和子开关轮询被限定为同一时间内只能存在一种，因此弃用这个接口和响应
	// case uint16(modbus.MSGTYPE_GETSNAPSHOOT_RESULT):
	// 	g.Debugln("主机快照轮询:")
	// 	if len(data16) <= 5 {
	// 		return
	// 	}
	// 	result := modbus.SnapShootMultiUploadDecode(data16[1:])
	// 	//转发给客户端
	// 	dl := bean.DeviceLoop{DeviceId: deviceId, DeviceCount: int(result.DeviceCount), LeftTime: int(result.LeftTime)}
	// 	var ssss []bean.SwitchSnapShoot
	// 	for _, v := range result.DeviceSnapShoots {
	// 		ssss = append(ssss, bean.SwitchSnapShoot{
	// 			Addr:         int(v.DeviceAddr),
	// 			DeviceType:   modbus.DEVICETYPE_MSG[v.DeviceType],
	// 			RatedCurrent: int(v.RatedCurrent),
	// 			PowerP:       v.Data.(modbus.DeviceSnapShootContent3).PowerP,
	// 			Energy_pt:    v.Data.(modbus.DeviceSnapShootContent3).Energy_pt,
	// 		})
	// 	}
	// 	dl.List = ssss
	// 	model.Broadcast(deviceId, func(ses cellnet.Session) {
	// 		if ses == nil {
	// 			return
	// 		}
	// 		ses.Send(&dl)
	// 	})
	// 	//轮询控制
	// 	if result.LeftTime == 0 {
	// 		model.SingleSend(deviceId, func(product string) {
	// 			g.Debugf("重开设备 %s 轮询\n", deviceId)
	// 			command.HostSnapShoot(deviceId, product, uint8(config.Conf.LoopInterval), 255)
	// 		}, config.Conf.CmdInterval)
	// 	}

	//开关实时数据轮询
	case uint16(modbus.MSGTYPE_RUNTIME_RESULT):
		g.Debugln("开关实时轮询:")
		if len(data16) <= 8 {
			return
		}
		result := modbus.ReadDeviceRuntimeUploadResult(data16)
		//转发给客户端
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.SwitchLoop{
				DeviceId:     deviceId,
				Addr:         int(result.Addr),
				DeviceType:   int(result.DeviceType),
				RatedCurrent: int(result.RatedCurrent),
				LeftTime:     int(result.LeftTime),
				Runtime: bean.SwitchRuntime1{
					HardFault:     int(result.Data.(modbus.DeviceRegisterDataRead1).HardFault),
					AlarmState:    result.Data.(modbus.DeviceRegisterDataRead1).AlarmState,
					ErrorState:    result.Data.(modbus.DeviceRegisterDataRead1).ErrorState,
					SwitchState:   int(result.Data.(modbus.DeviceRegisterDataRead1).SwitchState),
					SwitchMode:    int(result.Data.(modbus.DeviceRegisterDataRead1).SwitchMode),
					Rms_u:         result.Data.(modbus.DeviceRegisterDataRead1).Rms_U,
					Rms_i:         result.Data.(modbus.DeviceRegisterDataRead1).Rms_I,
					Power_p:       int(result.Data.(modbus.DeviceRegisterDataRead1).PowerP),
					Power_q:       int(result.Data.(modbus.DeviceRegisterDataRead1).PowerQ),
					Energy_pt:     result.Data.(modbus.DeviceRegisterDataRead1).Energy_pt,
					Energy_pt_l:   result.Data.(modbus.DeviceRegisterDataRead1).Energy_pt_l,
					Pf:            result.Data.(modbus.DeviceRegisterDataRead1).Pf,
					Freq:          result.Data.(modbus.DeviceRegisterDataRead1).Freq,
					Rms_IL:        int(result.Data.(modbus.DeviceRegisterDataRead1).Rms_il),
					Temperature:   result.Data.(modbus.DeviceRegisterDataRead1).Temperature,
					MotorRunTimes: int(result.Data.(modbus.DeviceRegisterDataRead1).MotorRunTimes),
					ErrorTimes:    int(result.Data.(modbus.DeviceRegisterDataRead1).ErrorTimes),
				},
			})
		})
		//轮询控制
		if result.LeftTime == 0 {
			model.SingleSend(deviceId, func(product string) {
				g.Debugf("重开设备 %s 开关 %d 轮询\n", deviceId, result.Addr)
				command.SwitchRuntime(deviceId, product, result.Addr, uint8(config.Conf.LoopInterval), 255)
			}, config.Conf.CmdInterval)
		}

	//开关设置数据上传（设备上电时主动上传，0xb3）
	//0xba，是0xb9请求开关设置数据的响应
	case uint16(modbus.MSGTYPE_UPLOADVALUES), uint16(modbus.MSGTYPE_SETTING_RESULT):
		g.Debugln("开关设置数据上传:")
		result := modbus.ReadDeviceSettingUploadResult(data16)
		ssu := bean.SwitchSettingUpload{DeviceId: deviceId, Product: product, Addr: int(result.Addr),
			DeviceType: int(result.DeviceType), RatedCurrent: int(result.RatedCurrent), Data: bean.SwitchSetting1{
				AlarmEnable:           result.Data.(modbus.DeviceRegisterDataWrite1).AlarmEnable,
				ErrorEnable:           result.Data.(modbus.DeviceRegisterDataWrite1).ErrorEnable,
				VoltageLimitRstEnable: result.Data.(modbus.DeviceRegisterDataWrite1).VoltageLimitRstEnable,
				IH_P:                  result.Data.(modbus.DeviceRegisterDataWrite1).IH_P,
				IH:                    result.Data.(modbus.DeviceRegisterDataWrite1).IH,
				UH_P:                  result.Data.(modbus.DeviceRegisterDataWrite1).UH_P,
				UH:                    result.Data.(modbus.DeviceRegisterDataWrite1).UH,
				UL_P:                  result.Data.(modbus.DeviceRegisterDataWrite1).UL_P,
				UL:                    result.Data.(modbus.DeviceRegisterDataWrite1).UL,
				PH_P:                  result.Data.(modbus.DeviceRegisterDataWrite1).PH_P,
				PH:                    result.Data.(modbus.DeviceRegisterDataWrite1).PH,
				EH_P:                  result.Data.(modbus.DeviceRegisterDataWrite1).EH_P,
				EH:                    result.Data.(modbus.DeviceRegisterDataWrite1).EH,
				IL_P:                  result.Data.(modbus.DeviceRegisterDataWrite1).IL_P,
				IL:                    result.Data.(modbus.DeviceRegisterDataWrite1).IL,
				TH_P:                  result.Data.(modbus.DeviceRegisterDataWrite1).TH_P,
				TH:                    result.Data.(modbus.DeviceRegisterDataWrite1).TH,
				UHL_CT:                result.Data.(modbus.DeviceRegisterDataWrite1).UHL_CT,
				UHL_RT:                result.Data.(modbus.DeviceRegisterDataWrite1).UHL_RT,
				IH_PH_CT:              result.Data.(modbus.DeviceRegisterDataWrite1).IH_PH_CT,
			}}
		var timers []bean.SwitchTimer1
		for _, v := range result.Timers.([]modbus.DeviceTimer1) {
			tmp := bean.SwitchTimer1{
				Group: v.Group,
				Type:  v.Type,
				Cycle: v.Cycle,
				State: v.State,
			}
			switch v.Cycle {
			case 0:
				tmp.Timer = bean.Once1{
					Timestamp: v.Timer.(modbus.Once1).Timestamp,
				}
			case 1:
				tmp.Timer = bean.Daily1{
					Hour:   v.Timer.(modbus.Daily1).Hour,
					Minute: v.Timer.(modbus.Daily1).Minute,
				}
			case 2:
				tmp.Timer = bean.Weekly1{
					WeekDay: v.Timer.(modbus.Weekly1).WeekDay,
					Hour:    v.Timer.(modbus.Weekly1).Hour,
					Minute:  v.Timer.(modbus.Weekly1).Minute,
				}
			case 3:
				tmp.Timer = bean.Monthly1{
					Day:    v.Timer.(modbus.Monthly1).Day,
					Hour:   v.Timer.(modbus.Monthly1).Hour,
					Minute: v.Timer.(modbus.Monthly1).Minute,
				}
			}
			timers = append(timers, tmp)
		}
		ssu.Timers = timers
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&ssu)
		})

	//开关状态事件
	case uint16(modbus.MSGTYPE_SWITCHSTATE):
		g.Debugln("开关状态事件:")
		result := modbus.SwitchStatusEventUploadResult(data16)
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.SwitchStateEvent{
				DeviceId:     deviceId,
				Addr:         int(result.Addr),
				DeviceType:   int(result.DeviceType),
				RatedCurrent: int(result.RatedCurrent),
				SwitchStatus: int(result.SwitchStatus),
			})
		})

	//开关模式事件
	case uint16(modbus.MSGTYPE_SWITCHMODE):
		g.Debugln("开关模式事件:")
		result := modbus.SwitchModeEventUploadResult(data16)
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.SwitchModeEvent{
				DeviceId:     deviceId,
				Addr:         int(result.Addr),
				DeviceType:   int(result.DeviceType),
				RatedCurrent: int(result.RatedCurrent),
				SwitchMode:   int(result.SwitchMode),
			})
		})

	//开关故障事件
	case uint16(modbus.MSGTYPE_ERROR):
		g.Debugln("开关故障事件:")
		result := modbus.DeviceFaultEventUploadResult(data16)
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.SwitchFaultEvent{
				DeviceId:     deviceId,
				Addr:         result.Addr,
				DeviceType:   result.DeviceType,
				RatedCurrent: result.RatedCurrent,
				FaultEvents:  result.FaultEvents,
				Data: bean.FaultEventData{
					OverCurrent:      result.Data.(modbus.FaultData1).OverCurrent,
					OverVoltage:      result.Data.(modbus.FaultData1).OverVoltage,
					UnderVoltage:     result.Data.(modbus.FaultData1).UnderVoltage,
					Overdrive:        result.Data.(modbus.FaultData1).Overdrive,
					ElectricQuantity: result.Data.(modbus.FaultData1).ElectricQuantity,
					OverHeat:         result.Data.(modbus.FaultData1).OverHeat,
					ShortOut:         result.Data.(modbus.FaultData1).ShortOut,
					ElectricLeakage:  result.Data.(modbus.FaultData1).ElectricLeakage,
					GroundElectrode:  result.Data.(modbus.FaultData1).GroundElectrode,
					Ack:              result.Data.(modbus.FaultData1).Ack,
				},
			})
		})

	//开关预警事件
	case uint16(modbus.MSGTYPE_ALARM):
		g.Debugln("开关预警事件:")
		result := modbus.DeviceWarnEventUploadResult(data16)
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.SwitchWarnEvent{
				DeviceId:     deviceId,
				Addr:         result.Addr,
				DeviceType:   result.DeviceType,
				RatedCurrent: result.RatedCurrent,
				WarnEvents:   result.WarnEvents,
				Data: bean.WarnEventData{
					OverCurrent:      result.Data.(modbus.WarnData1).OverCurrent,
					OverVoltage:      result.Data.(modbus.WarnData1).OverVoltage,
					UnderVoltage:     result.Data.(modbus.WarnData1).UnderVoltage,
					Overdrive:        result.Data.(modbus.WarnData1).Overdrive,
					ElectricQuantity: result.Data.(modbus.WarnData1).ElectricQuantity,
					OverHeat:         result.Data.(modbus.WarnData1).OverHeat,
					ShortOut:         result.Data.(modbus.WarnData1).ShortOut,
					ElectricLeakage:  result.Data.(modbus.WarnData1).ElectricLeakage,
					GroundElectrode:  result.Data.(modbus.WarnData1).GroundElectrode,
					Ack:              result.Data.(modbus.WarnData1).Ack,
				},
			})
		})

	//开关硬件故障事件
	case uint16(modbus.MSGTYPE_HARDWAREFAULT):
		g.Debugln("开关硬件故障事件:")
		result := modbus.DeviceHardwareFaultUploadResult(data16)
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.SwitchHardwareFaultEvent{
				DeviceId:       deviceId,
				Addr:           result.Addr,
				DeviceType:     result.DeviceType,
				RatedCurrent:   result.RatedCurrent,
				HardwareEvents: result.HardwareEvents,
			})
		})

	//开关漏电检测事件
	case uint16(modbus.MSGTYPE_ELECTRICLEAKAGETESTUPLOAD):
		g.Debugln("开关漏电检测事件:")
		result := modbus.ElectricLeakageTestEventUploadResult(data16)
		model.Broadcast(deviceId, func(ses cellnet.Session) {
			if ses == nil {
				return
			}
			ses.Send(&bean.SwitchElectricLeakageTestEvent{
				DeviceId:     deviceId,
				Addr:         result.Addr,
				DeviceType:   result.DeviceType,
				RatedCurrent: result.RatedCurrent,
				Status:       result.Status,
			})
		})

	//cmd
	case uint16(modbus.MSGTYPE_CMDRESULT):
		g.Debugln("modbus cmd响应:")
		switch data16[2] {
		case modbus.READRANGE_SUCCESS, modbus.READRANGE_FAILED:
			result := modbus.ReadRangeDecode(data16[1:])
			//主机多寄存器（6-10）
			if result.DeviceAddr == 0 {
				if result.Ok {
					modbus.ReadRangeSuccessDecodeWithPtr(result, data16[1:len(data16)-2], &modbus.HostRegisterData{})
					g.Debugf("主机多寄存器读成功，数据:%#v\n", result.Data)
					var deviceNum = int(result.Data.(*modbus.HostRegisterData).DeviceCount)
					//管理开关实时数据轮询
					model.UpdateDeviceCacheSwitchNumAndSend(deviceId, deviceNum, func(i int) {
						g.Debugf("关闭设备 %s 开关 %d 轮询\n", deviceId, i)
						command.SwitchRuntime(deviceId, product, uint8(i), uint8(config.Conf.LoopInterval), 0)
					}, config.Conf.CmdInterval)
					//loop更新设备子开关数，第一次赋值是就立即触发一次全子开关runtime读
					model.UpdateLoopCacheSwitchNumAndCall(deviceId, deviceNum, nil)
					var t1 []time.Duration
					var f1 []func()
					for i := 1; i <= deviceNum; i++ {
						t := i
						t1 = append(t1, config.Conf.CmdInterval)
						f1 = append(f1, func() {
							g.Debugf("单次请求设备 %s 子开关 %d runtime\n", deviceId, t)
							command.SwitchRuntimeOnce(deviceId, product, uint8(t))
						})
					}
					model.MultiSend(deviceId, f1, t1)
					//返回给客户端
					model.Broadcast(deviceId, func(ses cellnet.Session) {
						if ses == nil {
							return
						}
						ses.Send(&bean.DeviceAllRegister{
							DeviceId:        deviceId,
							SignalIntensity: int(result.Data.(*modbus.HostRegisterData).SignalIntensity),
							InternetMode:    int(result.Data.(*modbus.HostRegisterData).SignalIntensity),
							DeviceCount:     int(result.Data.(*modbus.HostRegisterData).DeviceCount),
							UploadInterval:  int(result.Data.(*modbus.HostRegisterData).UploadInterval),
						})
					})
				} else {
					g.Debugf("主机多寄存器读错误，错误原因:%s\n", result.ErrorMsg)
				}
				//子开关runtime
			} else if result.DataSize == 48 {
				if result.Ok {
					modbus.ReadRangeSuccessDecodeWithPtr(result, data16[1:len(data16)-2], &modbus.DeviceRegisterDataRead1{})
					g.Debugf("开关 %d runtime读成功，数据:%#v\n", result.DeviceAddr, result.Data)
					resp := &bean.SwitchRuntimeResponse{DeviceId: deviceId, Addr: int(result.DeviceAddr), Data: bean.SwitchRuntime1{
						HardFault:     int(result.Data.(*modbus.DeviceRegisterDataRead1).HardFault),
						AlarmState:    result.Data.(*modbus.DeviceRegisterDataRead1).AlarmState,
						ErrorState:    result.Data.(*modbus.DeviceRegisterDataRead1).ErrorState,
						SwitchState:   int(result.Data.(*modbus.DeviceRegisterDataRead1).SwitchState),
						SwitchMode:    int(result.Data.(*modbus.DeviceRegisterDataRead1).SwitchMode),
						Rms_u:         result.Data.(*modbus.DeviceRegisterDataRead1).Rms_U,
						Rms_i:         result.Data.(*modbus.DeviceRegisterDataRead1).Rms_I,
						Power_p:       int(result.Data.(*modbus.DeviceRegisterDataRead1).PowerP),
						Power_q:       int(result.Data.(*modbus.DeviceRegisterDataRead1).PowerQ),
						Energy_pt:     result.Data.(*modbus.DeviceRegisterDataRead1).Energy_pt,
						Energy_pt_l:   result.Data.(*modbus.DeviceRegisterDataRead1).Energy_pt_l,
						Pf:            result.Data.(*modbus.DeviceRegisterDataRead1).Pf,
						Freq:          result.Data.(*modbus.DeviceRegisterDataRead1).Freq,
						Rms_IL:        int(result.Data.(*modbus.DeviceRegisterDataRead1).Rms_il),
						Temperature:   result.Data.(*modbus.DeviceRegisterDataRead1).Temperature,
						MotorRunTimes: int(result.Data.(*modbus.DeviceRegisterDataRead1).MotorRunTimes),
						ErrorTimes:    int(result.Data.(*modbus.DeviceRegisterDataRead1).ErrorTimes),
					}}
					model.Broadcast(deviceId, func(ses cellnet.Session) {
						if ses == nil {
							return
						}
						ses.Send(resp)
					})
				} else {
					g.Debugf("开关 %d runtime读错误，错误原因:%s\n", result.DeviceAddr, result.ErrorMsg)
				}
			}

		//modbus单寄存器写响应
		case modbus.WRITESINGLE_SUCCESS, modbus.WRITESINGLE_FAILED:
			result := modbus.WriteSingleDecode(data16[1:])
			model.Broadcast(deviceId, func(ses cellnet.Session) {
				if ses == nil {
					return
				}
				ses.Send(&bean.ModbusSingleWriteResult{
					DeviceId:      deviceId,
					Product:       product,
					Ok:            result.Ok,
					Addr:          int(result.DeviceAddr),
					RegisterIndex: int(result.RegisterIndex),
					NewValue:      int(result.NewValue),
					ErrorMsg:      result.ErrorMsg,
				})
			})
		}
	}
}
