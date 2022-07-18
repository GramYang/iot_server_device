package handler

import (
	"iot_server4/bean"
	"iot_server4/cache"
	"iot_server4/command"
	"iot_server4/config"
	"iot_server4/model"
	"iot_server4/util"
	"time"

	sc "iot_server4/sqlx_client"

	g "github.com/GramYang/gylog"
	"github.com/davyxu/cellnet"
)

//处理心跳包的封装，所有的请求消息都带心跳包了。返回验证是否通过
func handleHeartbeat(msg *bean.ClientHeartBeat, ses cellnet.Session) bool {
	if msg.Token == "" || msg.Username == "" {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "username or token invalid",
		})
		ses.Close()
		return false
	}
	u := model.SessionToUser(ses)
	t := cache.GetToken(msg.Username)
	if u != nil { //非第一次收到心跳包
		//更新心跳
		u.LastPingTime = time.Now()
		//token过期了
		if t == "" {
			ses.Send(&bean.Shutdown{
				Code:    2,
				Message: "token expired",
			})
			ses.Close()
			return false
			//同用户名登录后将前者踢下线
		} else {
			if msg.Token != t {
				g.Debugf("同账号登录，关闭连接 %d msg.token %s cache token %s\n", ses.ID(), msg.Token, cache.GetToken(msg.Username))
				ses.Send(&bean.Shutdown{
					Code:    1,
					Message: "another client use this username to login",
				})
				ses.Close()
				return false
			}
		}
	} else { //第一次收到心跳包
		//先验证jwt
		if !util.ParseToken(msg.Token) {
			ses.Send(&bean.Shutdown{
				Code:    3,
				Message: "token error",
			})
			ses.Close()
			return false
		}
		//连接初始化
		_, err := model.CreateUser(ses)
		if err != nil {
			ses.Send(&bean.Shutdown{
				Code:    5,
				Message: "server error",
			})
			ses.Close()
			g.Errorln(ses.ID(), err)
			return false
		}
		//初次登录或者token失效 ||同用户名登录后将前者踢下线
		if t == "" || msg.Token != t {
			//存入/覆盖token，刷新token生命周期
			cache.SaveToken(msg.Username, msg.Token)
		}
	}
	return true
}

//每次消息发给客户端都会收到一个心跳包，因此需要取消间隔过近的延迟响应
func handler1(msg *bean.ClientHeartBeat, ses cellnet.Session) {
	ok := handleHeartbeat(msg, ses)
	if !ok {
		return
	}
	//延迟发送心跳包响应，且控制响应的间隔时间为心跳时间
	u := model.SessionToUser(ses)
	if u.TimerHandler != nil {
		u.TimerHandler.Stop()
	}
	u.TimerHandler = time.AfterFunc(time.Second*config.Conf.HeartbeatInterval, func() {
		ses.Send(&bean.ClientHeartBeatResponse{
			Username: msg.Username,
			Token:    msg.Token,
		})
	})
}

func handler2(msg *bean.GetDeviceList, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	deviceList, err := sc.GetDeviceListByUsername(msg.Heartbeat.Username)
	if err != nil {
		ses.Send(&bean.ResultMessage{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	} else {
		if len(deviceList) == 0 {
			return
		}
		var deviceIds, products []string
		for _, v := range deviceList {
			deviceIds = append(deviceIds, v.DeviceId)
			products = append(products, v.Product)

		}
		//保存客户端绑定设备列表
		model.SaveDeviceList(ses.ID(), deviceIds)
		//添加devicecache，并发送初始化命令
		model.AddDeviceCachesAndSend(deviceIds, products, ses.ID(), nil, config.Conf.CmdInterval)
		for _, v := range deviceIds {
			model.SingleSend(v, func(product string) {
				g.Debugf("主机 %s 第一个全寄存器读\n", v)
				command.HostAllRegister(v, product)
			}, config.Conf.CmdInterval)
		}
		//添加devicesession，因为命令已经发出
		model.AddDeviceSessions(deviceIds, time.Now().UnixMilli())
		//添加loopCache
		model.AddLoopCaches(deviceIds, products)
		//用户绑定设备列表返回给客户端
		ses.Send(&bean.GetDeviceListResult{
			DeviceIds: deviceIds,
			Products:  products,
		})
	}
}

func handler38(msg *bean.ClientLogout, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	cache.DeleteToken(msg.Heartbeat.Username)
}
