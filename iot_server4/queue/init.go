package queue

var CacheQueue Queue

func SetUp(r float64) {
	CacheQueue = NewQueue()
	CacheQueue.StartLoopWithRate(r)
}
