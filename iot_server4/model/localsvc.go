package model

import (
	"sync"
	"time"

	"github.com/davyxu/cellnet"
)

var (
	localServices      []cellnet.Peer
	localServicesGuard sync.RWMutex
)

func AddLocalService(p cellnet.Peer) {
	localServicesGuard.Lock()
	localServices = append(localServices, p)
	localServicesGuard.Unlock()
}

func VisitLocalService(callback func(cellnet.Peer) bool) {
	localServicesGuard.RLock()
	defer localServicesGuard.RUnlock()

	for _, svc := range localServices {
		if !callback(svc) {
			break
		}
	}
}

func IsAllReady() (ret bool) {
	ret = true
	VisitLocalService(func(svc cellnet.Peer) bool {
		if !svc.(cellnet.PeerReadyChecker).IsReady() {
			ret = false
			return false
		}

		return true
	})

	return
}

func CheckReady() {
	for {
		time.Sleep(time.Second * 3)
		if IsAllReady() {
			break
		}
	}
}

func StopAllService() {
	localServicesGuard.RLock()
	defer localServicesGuard.RUnlock()

	for i := len(localServices) - 1; i >= 0; i-- {
		svc := localServices[i]
		svc.Stop()
	}
}
