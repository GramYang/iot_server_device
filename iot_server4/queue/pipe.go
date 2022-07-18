package queue

import "sync"

type pipe struct {
	list      []interface{}
	listGuard sync.Mutex
	listCond  *sync.Cond
}

func (p *pipe) Add(msg interface{}) {
	p.listGuard.Lock()
	p.list = append(p.list, msg)
	p.listGuard.Unlock()
	p.listCond.Signal()
}

func (p *pipe) Count() int {
	p.listGuard.Lock()
	defer p.listGuard.Unlock()
	return len(p.list)
}

func (p *pipe) Reset() {
	p.listGuard.Lock()
	p.list = p.list[0:0]
	p.listGuard.Unlock()
}

func (p *pipe) Pick(list *[]interface{}) (exit bool) {
	//在Pipe的list长度为0时，阻塞。上面调用了Add向Pipe的list添加元素后解除阻塞
	p.listGuard.Lock()
	for len(p.list) == 0 {
		//fmt.Println(time.Now().UnixNano()/1000000)//在pipe为空阻塞时输出时间，debug用
		p.listCond.Wait()
	}
	p.listGuard.Unlock()
	//将Pipe的list中的值复制进参数list，然后清空Pipe的list
	p.listGuard.Lock()
	for _, data := range p.list {
		if data == nil {
			exit = true //在list里面有nil时才会退出循环返回true
			break
		} else {
			*list = append(*list, data)
		}
	}
	p.list = p.list[0:0]
	p.listGuard.Unlock()
	return
}

func newPipe() *pipe {
	self := &pipe{}
	self.listCond = sync.NewCond(&self.listGuard)
	return self
}
