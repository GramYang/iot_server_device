package handler

import (
	"iot_server4/bean"
	"iot_server4/command"
	"iot_server4/config"
	"iot_server4/model"
	"iot_server4/prom"
	"time"

	g "github.com/GramYang/gylog"
	"github.com/davyxu/cellnet"

	q "iot_server4/queue"
)

func Handler2(ev cellnet.Event) {
	switch msg := ev.Message().(type) {

	//打印开启连接信息
	case *cellnet.SessionAccepted:
		g.Debugf("开启连接 %d\n", ev.Session().ID())
		prom.ClientOnline()

	//连接关闭，此消息在客户端关闭连接和后端心跳检查关闭连接时触发
	case *cellnet.SessionClosed:
		prom.ClientOffline()
		model.DeleteCacheSession(ev.Session().ID(), model.QueryDeviceList(ev.Session().ID()), func(name, product string, switchCount int) {
			//关闭该设备的主机和开关轮询，这里deviceCache里的吸管已经关闭了，直接用495消息队列发
			// time.AfterFunc(config.Conf.CmdInterval*time.Millisecond, func() {
			// 	q.CacheQueue.Post(func() {
			// 		g.Debugf("关闭设备 %s 轮询\n", name)
			// 		command.HostSnapShoot(name, product, uint8(config.Conf.LoopInterval), 0)
			// 	})
			// })
			for i := 1; i <= switchCount; i++ {
				t := i
				time.AfterFunc(config.Conf.CmdInterval*time.Millisecond, func() {
					q.CacheQueue.Post(func() {
						g.Debugf("关闭设备 %s 开关 %d 轮询\n", name, t)
						command.SwitchRuntime(name, product, uint8(t), uint8(config.Conf.LoopInterval), 0)
					})
				})
			}
		})
		//loopCache移除设备
		model.DeleteLoopCache(model.QueryDeviceList(ev.Session().ID()))

	//测试消息用于测试websocket连接，直接发回去
	case *bean.TestEchoACK:
		ev.Session().Send(msg)

	//无响应，触发逻辑
	case *bean.ClientHeartBeat:
		handler1(msg, ev.Session())

	//响应+触发逻辑
	case *bean.GetDeviceList:
		handler2(msg, ev.Session())

	//响应开关事件
	case *bean.MultiSwitchOperation:
		handler3(msg, ev.Session())

	//响应开关事件
	case *bean.SwitchOn:
		handler4(msg, ev.Session())

	//响应开关事件
	case *bean.SwitchOff:
		handler5(msg, ev.Session())

	//响应开关事件
	case *bean.SwitchLockOn:
		handler6(msg, ev.Session())

	//响应开关事件
	case *bean.SwitchLockOff:
		handler7(msg, ev.Session())

	//响应：单条写成功，设备编号:0 寄存器:0 最新值:5
	//还有一堆设备上电后主动发送的消息，包括：
	//主机全寄存器上报、主机下全设备上报、所有子开关的开关状态事件、所有子开关的设置数据上报
	case *bean.DeviceReset:
		handler8(msg, ev.Session())

	//所有对设置数据的修改都会触发一次“设置数据上报”，即0xb3
	case *bean.SetTimer:
		handler9(msg, ev.Session())

	//响应“漏电检测事件上报”
	case *bean.SwitchElectricLeakageTest:
		handler10(msg, ev.Session())

	//响应“设置数据上报”
	//0过流1过压2欠压3过载4电量5过温
	case *bean.SwitchAlarmEnable:
		handler11(msg, ev.Session())

	//响应“设置数据上报”
	//0过流1过压2欠压3过载4电量5过温6电弧
	case *bean.SwitchErrorEnable:
		handler12(msg, ev.Session())

	//响应
	case *bean.DeviceElectricQuantity:
		handler13(msg, ev.Session())

	//请求主机下所有设备状态
	case *bean.DeviceAllState:
		handler14(msg, ev.Session())

	//请求开关设置数据
	case *bean.GetSwitchSetting:
		handler15(msg, ev.Session())

	//开启子开关轮询
	case *bean.SwitchLoopOn:
		handler16(msg, ev.Session())

	//关闭子开关轮询
	case *bean.SwitchLoopOff:
		handler17(msg, ev.Session())

	//子开关清除故障
	case *bean.SwitchClearFault:
		handler18(msg, ev.Session())

	//过欠压恢复设置
	case *bean.VoltageLimitRstEnable:
		handler19(msg, ev.Session())

	case *bean.SetIHP:
		handler20(msg, ev.Session())

	case *bean.SetIH:
		handler21(msg, ev.Session())

	case *bean.SetUHP:
		handler22(msg, ev.Session())

	case *bean.SetUH:
		handler23(msg, ev.Session())

	case *bean.SetULP:
		handler24(msg, ev.Session())

	case *bean.SetUL:
		handler25(msg, ev.Session())

	case *bean.SetPHP:
		handler26(msg, ev.Session())

	case *bean.SetPH:
		handler27(msg, ev.Session())

	case *bean.SetEHP:
		handler28(msg, ev.Session())

	// case *bean.SetEH:
	// 	handler29(msg, ev.Session())

	case *bean.SetILP:
		handler30(msg, ev.Session())

	case *bean.SetIL:
		handler31(msg, ev.Session())

	case *bean.SetTHP:
		handler32(msg, ev.Session())

	case *bean.SetTH:
		handler33(msg, ev.Session())

	case *bean.SetUHLCT:
		handler34(msg, ev.Session())

	case *bean.SetUHLRT:
		handler35(msg, ev.Session())

	case *bean.SetIHPHCT:
		handler36(msg, ev.Session())

	//废弃了，这里转由后端控制
	// case *bean.GetAllSwitchRuntime:
	// 	handler37(msg, ev.Session())

	//客户端退出登录
	case *bean.ClientLogout:
		handler38(msg, ev.Session())
	}

}
