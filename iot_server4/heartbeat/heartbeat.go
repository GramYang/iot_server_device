package heartbeat

import (
	"iot_server4/command"
	"iot_server4/config"
	"iot_server4/model"
	"time"

	g "github.com/GramYang/gylog"
	"github.com/davyxu/cellnet/timer"
)

func StartCheck() {
	durationClient := config.Conf.HeartbeatInterval + 2
	durationDevice := config.Conf.HeartbeatInterval + 2
	//客户端心跳检查
	timeOutDurClient := time.Duration(durationClient) * time.Second
	g.Infof("client heatbeat duration: '%ds'\n", durationClient)
	timer.NewLoop(nil, timeOutDurClient, func(loop *timer.Loop) {
		now := time.Now()
		model.VisitUser(func(u *model.User) bool {
			if now.Sub(u.LastPingTime) > timeOutDurClient {
				g.Warningf("Close client due to heatbeat time out, id: %d\n", u.ClientSession.ID())
				//关闭session
				u.ClientSession.Close()
			}
			//这里对应的是sync.Map的Range接口，返回false则停止遍历
			return true
		})
	}, nil).Start()
	//设备心跳检查
	timeOutDurDevice := time.Duration(durationDevice) * time.Second
	g.Infof("device heatbeat duration: '%d's\n", durationDevice)
	timer.NewLoop(nil, timeOutDurDevice, func(loop *timer.Loop) {
		//轮询devicesession，下线时持续发送iccid心跳包
		model.VisitDeviceSession(time.Now().UnixMilli(), func(deviceId string) {
			model.SingleSend(deviceId, func(product string) {
				g.Debugf("设备 %s 发送iccid心跳\n", deviceId)
				command.Iccid(deviceId, product)
			}, config.Conf.CmdInterval)
		}, func(deviceId string) {
			//这个时候已经没有响应了，你发什么设备也收不到，所以什么都不发。
			//但是为了触发主机全寄存器读的子开关轮询接口，将该设备的子开关数置0，同时设备下线。
			model.UpdateDeviceCacheSwitchNum(deviceId, 0)
		})
	}, nil).Start()
	//主机全子开关runtime轮询请求
	SwitchRuntimeDur := timeOutDurDevice * 2
	timer.NewLoop(nil, SwitchRuntimeDur, func(loop *timer.Loop) {
		model.VisitLoopCache(func(deviceId, product string, num int) {
			var t1 []time.Duration
			var f1 []func()
			for i := 1; i <= num; i++ {
				t := i
				t1 = append(t1, config.Conf.CmdInterval)
				f1 = append(f1, func() {
					g.Debugf("单次请求设备 %s 子开关 %d runtime\n", deviceId, t)
					command.SwitchRuntimeOnce(deviceId, product, uint8(t))
				})
			}
			model.MultiSend(deviceId, f1, t1)
		})
	}, nil).Start()
}
