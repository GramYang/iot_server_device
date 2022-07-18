package queue

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Queue interface {
	StartLoop()
	StartLoopWithRate(r float64)
	StopLoop()
	Wait()
	Post(f func())
	EnableCapturePanic(v bool)
	Count() int
}

type CapturePanicNotifyFunc func(interface{}, Queue)

type queue struct {
	*pipe
	endSignal    sync.WaitGroup
	capturePanic bool
	onPanic      CapturePanicNotifyFunc
}

func (q *queue) StartLoop() {
	q.endSignal.Add(1)
	go func() {
		var writeList []interface{}
		for {
			writeList = writeList[0:0]
			exit := q.Pick(&writeList)
			for _, msg := range writeList {
				switch t := msg.(type) {
				case func():
					q.protectedCall(t)
				case nil:
					break
				default:
					fmt.Printf("unexpected type %T\n", t)
				}
			}
			if exit {
				break
			}
		}
		q.endSignal.Done()
	}()
}

func (q *queue) StartLoopWithRate(r float64) {
	q.endSignal.Add(1)
	go func() {
		var writeList []interface{}
		l := rate.NewLimiter(rate.Limit(r), 1)
		c, cancel := context.WithCancel(context.TODO())
		for {
			writeList = writeList[0:0]
			exit := q.Pick(&writeList)
			for _, msg := range writeList {
				switch t := msg.(type) {
				case func():
					err := l.Wait(c)
					if err != nil {
						cancel()
						return
					}
					q.protectedCall(t)
				case nil:
					break
				default:
					fmt.Printf("unexpected type %T\n", t)
				}
			}
			if exit {
				break
			}
		}
		cancel()
		q.endSignal.Done()
	}()
}

func (q *queue) StopLoop() {
	q.Add(nil)
}

func (q *queue) Wait() {
	q.endSignal.Wait()
}

func (q *queue) Post(callback func()) {
	if callback == nil {
		return
	}
	q.Add(callback)
}

func (q *queue) EnableCapturePanic(v bool) {
	q.capturePanic = v
}

func (q *queue) protectedCall(callback func()) {
	if q.capturePanic {
		defer func() {
			if err := recover(); err != nil {
				q.onPanic(err, q)
			}
		}()
	}
	callback()
}

func NewQueue() Queue {
	return &queue{
		pipe: newPipe(),
		onPanic: func(i interface{}, queue Queue) {
			fmt.Printf("%s: %v \n%s\n", time.Now().Format("2006-01-02 15:04:05"), i, string(debug.Stack()))
			debug.PrintStack()
		},
	}
}
