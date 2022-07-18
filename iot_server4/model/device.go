package model

import (
	"iot_server4/config"
	"sync"
	"time"

	q "iot_server4/queue"

	"github.com/davyxu/cellnet"

	g "github.com/GramYang/gylog"
)

var deviceCache map[string]*cache = make(map[string]*cache)
var rwmuxDevice sync.RWMutex

type cache struct {
	LastCMDTime  int64  //最近的客户端命令毫秒时间，记录的是开始发送时间而不是发送完成时间
	Product      string //加上这个用于modbus命令发送
	SwitchNum    int    //设备下开关数量，用于设备实时轮询
	SessionIds   []int64
	Straw        q.Queue //客户端转发+服务端发送的设备指令都通过此容器发送，以控制延迟
	SwitchLoopOn []int   //标记子开关轮询开启，下标表示子开关地址，0表示开启，1或者1以上表示关闭
}

//添加多个设备缓存，初次添加的设备发送多个初始化命令，不需要判断设备是否在线
func AddDeviceCachesAndSend(deviceNames, products []string, sesId int64, multiCall []func(string, string), delay time.Duration) {
	rwmuxDevice.Lock()
	defer rwmuxDevice.Unlock()
	for i := 0; i < len(deviceNames); i++ {
		t1 := i
		c, ok := deviceCache[deviceNames[t1]]
		if ok {
			c.SessionIds = append(c.SessionIds, sesId)
		} else {
			straw := q.NewQueue()
			straw.StartLoop()
			deviceCache[deviceNames[t1]] = &cache{Product: products[t1], SessionIds: []int64{sesId}, Straw: straw}
			if len(multiCall) > 0 {
				for j := 0; j < len(multiCall); j++ {
					t2 := j
					straw.Post(func() {
						time.AfterFunc(delay*time.Millisecond, func() {
							q.CacheQueue.Post(func() { multiCall[t2](deviceNames[t1], products[t1]) })
						})
					})
				}
			}
		}
	}
}

//检查设备状态（是否注册、是否在线、上次命令是否在config.Conf.CmdInterval*3毫秒内），正常则更新LastCMDTime
//这个接口是转发客户端设备指令时使用的，因为存在多个客户端并发发送指令，需要过滤，外加一个Straw回调
func CheckDeviceCacheAndSessionAndSend(deviceName string, now int64, f func(), delay time.Duration) (bool, string) {
	rwmuxDevice.RLock()
	defer rwmuxDevice.RUnlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return false, "device is not initialized"
	}
	if !CheckDeviceSessionOnline(deviceName) {
		return false, "device offline"
	}
	g.Debugf("设备 %s 命令延迟毫秒数 %d 延迟判断间隔 %d\n", deviceName, now-c.LastCMDTime, config.Conf.CmdInterval*2)
	if now-c.LastCMDTime <= int64(config.Conf.CmdInterval*2) {
		return false, "other cmds on operating"
	}
	c.LastCMDTime = now
	c.Straw.Post(func() {
		time.AfterFunc(delay*time.Millisecond, func() {
			q.CacheQueue.Post(func() { f() })
		})
	})
	return true, ""
}

//上面接口专门为延迟发送子开关命令的定制版
func CheckDeviceCacheAndSessionAndSend1(deviceName string, now int64, f func(int), delay time.Duration) (bool, string) {
	rwmuxDevice.RLock()
	defer rwmuxDevice.RUnlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return false, "device is not initialized"
	}
	if !CheckDeviceSessionOnline(deviceName) {
		return false, "device offline"
	}
	g.Debugf("设备 %s 命令延迟毫秒数 %d 延迟判断间隔 %d\n", deviceName, now-c.LastCMDTime, config.Conf.CmdInterval*2)
	if now-c.LastCMDTime <= int64(config.Conf.CmdInterval*2) {
		return false, "other cmds on operating"
	}
	c.LastCMDTime = now
	for i := 1; i <= c.SwitchNum; i++ {
		t := i
		c.Straw.Post(func() {
			time.AfterFunc(delay*time.Millisecond, func() {
				q.CacheQueue.Post(func() { f(t) })
			})
		})
	}
	return true, ""
}

//遍历devicecache中的deviceList移除指定session，cb是当SessionIds为空移除entry时的回调
func DeleteCacheSession(sesId int64, deviceList []string, cb func(string, string, int)) {
	rwmuxDevice.Lock()
	defer rwmuxDevice.Unlock()
	if len(deviceCache) == 0 {
		return
	}
	if len(deviceList) == 0 {
		return
	}
	for _, v := range deviceList {
		c, ok := deviceCache[v]
		if !ok {
			continue
		}
		var tmp []int64
		for _, v1 := range c.SessionIds {
			if v1 != sesId {
				tmp = append(tmp, v1)
			}

		}
		c.SessionIds = tmp
		if len(c.SessionIds) == 0 {
			c.Straw.StopLoop() //关闭
			cb(v, c.Product, c.SwitchNum)
			DeleteDeviceSession(v)
			delete(deviceCache, v)
			continue
		}
	}
}

//更新设备开关数目顺便管理开关实时数据的轮询
//更新设备中的子开关数量，并设置子开关轮询。
//两种情况触发子开关轮询：子开关数从0到指定值（初始化），子开关数变更（有变更）
func UpdateDeviceCacheSwitchNumAndSend(deviceName string, switchNum int, close func(index int), delay time.Duration) {
	rwmuxDevice.Lock()
	defer rwmuxDevice.Unlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return
	}
	if c.SwitchNum != switchNum {
		var oldSwitchNum = c.SwitchNum
		c.SwitchNum = switchNum
		if oldSwitchNum > 0 { //有就关闭老的开关轮询
			for i := 1; i <= oldSwitchNum; i++ {
				t := i
				c.Straw.Post(func() {
					time.AfterFunc(delay*time.Millisecond, func() {
						q.CacheQueue.Post(func() { close(t) })
					})
				})
			}
		}
		//为什么要在这里初始化SwitchLoopOn？因为你必须知道子开关数量后才能初始化
		var arr []int
		for i := 1; i <= switchNum; i++ {
			arr = append(arr, 0)
		}
		c.SwitchLoopOn = arr
	}
}

//设置switchnum
func UpdateDeviceCacheSwitchNum(deviceName string, switchNum int) {
	rwmuxDevice.Lock()
	defer rwmuxDevice.Unlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return
	}
	c.SwitchNum = switchNum
	var arr []int
	for i := 1; i <= switchNum; i++ {
		arr = append(arr, 0)
	}
	c.SwitchLoopOn = arr
}

//更新子开关轮询开启状态
func UpdateDeviceCacheSwitchLoopAndSend(deviceName string, addr int, isOpen bool, cb func(string), delay time.Duration) {
	rwmuxDevice.Lock()
	defer rwmuxDevice.Unlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return
	}
	if isOpen {
		if c.SwitchLoopOn[addr] == 0 {
			c.Straw.Post(func() {
				time.AfterFunc(delay*time.Millisecond, func() {
					q.CacheQueue.Post(func() { cb(c.Product) })
				})
			})
		}
		c.SwitchLoopOn[addr]++
	} else {
		c.SwitchLoopOn[addr]--
		if c.SwitchLoopOn[addr] == 0 {
			c.Straw.Post(func() {
				time.AfterFunc(delay*time.Millisecond, func() {
					q.CacheQueue.Post(func() { cb(c.Product) })
				})
			})
		}
	}
}

//遍历设备的SessionIds，调用cb
func Broadcast(deviceName string, cb func(cellnet.Session)) {
	rwmuxDevice.RLock()
	defer rwmuxDevice.RUnlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return
	} else if len(c.SessionIds) > 0 {
		for _, v := range c.SessionIds {
			cb(FrontendSessionManager.GetSession(v))
		}
	}
}

//单独发，带modbus用的product
func SingleSend(deviceName string, f func(string), delay time.Duration) {
	rwmuxDevice.RLock()
	defer rwmuxDevice.RUnlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return
	}
	c.Straw.Post(func() {
		time.AfterFunc(delay*time.Millisecond, func() {
			q.CacheQueue.Post(func() { f(c.Product) })
		})
	})
}

//发送多条命令，分别指定延迟
func MultiSend(deviceName string, fs []func(), delays []time.Duration) {
	rwmuxDevice.RLock()
	defer rwmuxDevice.RUnlock()
	c, ok := deviceCache[deviceName]
	if !ok {
		return
	}
	if len(fs) != len(delays) {
		return
	}
	for i := 0; i < len(fs); i++ {
		t := i
		c.Straw.Post(func() {
			time.AfterFunc(delays[t]*time.Millisecond, func() {
				q.CacheQueue.Post(func() { fs[t]() })
			})
		})
	}
}
