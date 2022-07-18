package model

import (
	"sync"
)

//此缓存专门用于全设备子开关轮询，因为要定期轮询，为了减轻device的压力单独列出来，因此和device同步
var loopCache map[string]*loop = make(map[string]*loop)
var rmwuxLoop sync.RWMutex

type loop struct {
	Product    string
	SwitchNum  int
	SessionNum int //这里用数字来指代关联客户端数
}

//添加设备
func AddLoopCaches(deviceIds, products []string) {
	rmwuxLoop.Lock()
	defer rmwuxLoop.Unlock()
	for i := 0; i < len(deviceIds); i++ {
		l, ok := loopCache[deviceIds[i]]
		if ok {
			l.SessionNum++
		} else {
			loopCache[deviceIds[i]] = &loop{Product: products[i], SessionNum: 1}
		}
	}
}

//移除设备
func DeleteLoopCache(deviceList []string) {
	rmwuxLoop.Lock()
	defer rmwuxLoop.Unlock()
	if len(loopCache) == 0 {
		return
	}
	if len(deviceList) == 0 {
		return
	}
	for _, v := range deviceList {
		l, ok := loopCache[v]
		if !ok {
			continue
		}
		l.SessionNum--
		if l.SessionNum == 0 {
			delete(loopCache, v)
		}
	}
}

//更新设备子开关数量
func UpdateLoopCacheSwitchNumAndCall(deviceId string, num int, callback func()) {
	rmwuxLoop.Lock()
	defer rmwuxLoop.Unlock()
	l, ok := loopCache[deviceId]
	if !ok {
		return
	}
	if l.SwitchNum == 0 && callback != nil {
		callback()
	}
	l.SwitchNum = num
}

//遍历loopCache执行子开关runtime查询
func VisitLoopCache(f func(string, string, int)) {
	rmwuxLoop.RLock()
	defer rmwuxLoop.RUnlock()
	if len(loopCache) == 0 {
		return
	}
	for deviceId, l := range loopCache {
		f(deviceId, l.Product, l.SwitchNum)
	}
}
