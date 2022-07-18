package model

import (
	"iot_server4/config"
	"iot_server4/prom"
	"sync"
)

var deviceSession map[string]*deviceStatus = make(map[string]*deviceStatus)
var rwmuxSession sync.RWMutex

type deviceStatus struct {
	DeviceOnline           bool  //客户端上线且在设备响应后标记此flag
	LastDeviceResponseTime int64 //毫秒数
}

//遍历deviceSession执行callback，更新设备在线状态
func VisitDeviceSession(now int64, heartbeat func(string), offline func(string)) {
	rwmuxSession.Lock()
	defer rwmuxSession.Unlock()
	if len(deviceSession) == 0 { //为空直接退出
		return
	}
	for deviceName, ses := range deviceSession {
		heartbeat(deviceName)
		ok := checkDeviceStatus(ses, now)
		if ok {
			ses.DeviceOnline = false
			offline(deviceName)
			prom.DeviceOffline()
		} else {
			ses.DeviceOnline = true
			prom.DeviceOnline()
		}
	}
}

//检查设备状态，true表示心跳超时
func checkDeviceStatus(d *deviceStatus, now int64) bool {
	if now-d.LastDeviceResponseTime > (int64(config.Conf.HeartbeatInterval+2) * 1000) {
		return true
	} else {
		return false
	}
}

//添加对应客户端的设备，跳过已存在的设备
func AddDeviceSession(deviceName string) {
	rwmuxSession.Lock()
	defer rwmuxSession.Unlock()
	if _, ok := deviceSession[deviceName]; ok {
		return
	} else {
		deviceSession[deviceName] = &deviceStatus{}
	}
}

//添加多个设备，传入第一个命令的时间来开启心跳判断
func AddDeviceSessions(deviceNames []string, resquestTime int64) {
	rwmuxSession.Lock()
	defer rwmuxSession.Unlock()
	for _, v := range deviceNames {
		if _, ok := deviceSession[v]; ok {
			continue
		} else {
			deviceSession[v] = &deviceStatus{LastDeviceResponseTime: resquestTime}
		}
	}
}

//更新设备信息并调用上线回调
func UpdateDeviceSession(deviceName string, responseTime int64, onlineCallback func(string)) {
	rwmuxSession.Lock()
	defer rwmuxSession.Unlock()
	if c, ok := deviceSession[deviceName]; ok {
		c.LastDeviceResponseTime = responseTime
		//单纯的更新响应时间
		if !c.DeviceOnline { //设备上线
			c.DeviceOnline = true
			onlineCallback(deviceName)
		}
	} else {
		return
	}
}

//检查设备是否在线，这也是session的核心功能，在设备下线的情况下返回消息给客户端
func CheckDeviceSessionOnline(deviceName string) bool {
	rwmuxSession.RLock()
	defer rwmuxSession.RUnlock()
	if c, ok := deviceSession[deviceName]; ok {
		return c.DeviceOnline
	} else {
		return false
	}
}

//移除设备
func DeleteDeviceSession(deviceName string) {
	rwmuxSession.Lock()
	defer rwmuxSession.Unlock()
	if _, ok := deviceSession[deviceName]; !ok {
		return
	}
	delete(deviceSession, deviceName)
}
