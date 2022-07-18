package model

import "sync"

var clientBindDevice map[int64][]string = make(map[int64][]string)
var rwmuxClient sync.RWMutex

func QueryDeviceList(sesId int64) []string {
	rwmuxClient.RLock()
	defer rwmuxClient.RUnlock()
	return clientBindDevice[sesId]
}

func SaveDeviceList(sesId int64, deviceList []string) {
	rwmuxClient.Lock()
	defer rwmuxClient.Unlock()
	clientBindDevice[sesId] = deviceList
}
