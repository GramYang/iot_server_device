package handler

import (
	"iot_server4/bean"
	"iot_server4/command"
	"iot_server4/config"
	"iot_server4/model"
	"time"

	g "github.com/GramYang/gylog"
	"github.com/davyxu/cellnet"
)

func handler3(msg *bean.MultiSwitchOperation, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Cmd == nil || len(msg.Cmd) == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "multi switch operation invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("连接 %d 一键开合闸 时间 %d\n", ses.ID(), time.Now().Unix())
	//一键分合闸在发送指令后设备会发出开关状态事件，有几个开关就有几个事件
	//延迟为开关数*200毫秒
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.HostSwitch(msg.DeviceId, msg.Product, msg.Cmd)
	}, config.Conf.CmdInterval*time.Duration(len(msg.Cmd)))
	if !ok {
		g.Debugf("连接 %d 结果消息 %s\n", ses.ID(), msg1)
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler4(msg *bean.SwitchOn, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch on invalid",
		})
		ses.Close()
		return
	}
	//开闸
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.SwitchOn(msg.DeviceId, msg.Product, uint8(msg.Addr))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler5(msg *bean.SwitchOff, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch off invalid",
		})
		ses.Close()
		return
	}
	//关闸
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.SwitchOff(msg.DeviceId, msg.Product, uint8(msg.Addr))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler6(msg *bean.SwitchLockOn, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch lock on invalid",
		})
		ses.Close()
		return
	}
	//锁定
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.SwitchLockOn(msg.DeviceId, msg.Product, uint8(msg.Addr))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler7(msg *bean.SwitchLockOff, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch lock off invalid",
		})
		ses.Close()
		return
	}
	//解锁
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.SwitchLockOff(msg.DeviceId, msg.Product, uint8(msg.Addr))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler8(msg *bean.DeviceReset, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "device reset invalid",
		})
		ses.Close()
		return
	}
	//重分配地址并锁定
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.HostResetAndLock(msg.DeviceId, msg.Product)
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler14(msg *bean.DeviceAllState, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "device all state invalid",
		})
		ses.Close()
		return
	}
	//读取主机下所有开关状态
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.HostAllSwitchState(msg.DeviceId, msg.Product)
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

// func handler37(msg *bean.GetAllSwitchRuntime, ses cellnet.Session) {
// 	ok := handleHeartbeat(&msg.Heartbeat, ses)
// 	if !ok {
// 		return
// 	}
// 	if msg.DeviceId == "" || msg.Product == "" {
// 		ses.Send(&bean.Shutdown{
// 			Code:    4,
// 			Message: "get all switch runtime invalid",
// 		})
// 		ses.Close()
// 		return
// 	}
// 	//请求主机下所有子开关的runtime
// 	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend1(msg.DeviceId, time.Now().UnixMilli(), func(addr int) {
// 		command.SwitchRuntimeOnce(msg.DeviceId, msg.Product, uint8(addr))
// 	}, config.Conf.CmdInterval)
// 	if !ok {
// 		ses.Send(&bean.DeviceResultMessage{
// 			IsSuccess: ok,
// 			DeviceId:  msg.DeviceId,
// 			Message:   msg1,
// 		})
// 	}
// }
