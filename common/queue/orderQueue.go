package queue

import "sync"

type OrderQueue struct {
	lock  sync.Mutex
	nodes []interface{}
	cmp   func(v1, v2 interface{}) int
}

func NewOrderQueue(cmp func(v1, v2 interface{}) int) *OrderQueue {
	return &OrderQueue{cmp: cmp}
}

func (oq *OrderQueue) EnQueue(node interface{}) {
	oq.lock.Lock()
	defer oq.lock.Unlock()
	var t []interface{}
	flag := false
	for i := 0; i < len(oq.nodes); i++ {
		v := oq.nodes[i]
		if !flag && oq.cmp(v, node) < 0 {
			flag = true
			t = append(t, node, oq.nodes[i])

		} else {
			t = append(t, oq.nodes[i])
		}
	}
	if !flag {
		t = append(t, node)
	}

	oq.nodes = t
}

func (oq *OrderQueue) DeQueue() interface{} {
	oq.lock.Lock()
	defer oq.lock.Unlock()

	l := len(oq.nodes)
	if l == 0 {
		return nil
	}

	v := oq.nodes[l-1]
	oq.nodes = oq.nodes[:l-1]
	return v
}

func (oq *OrderQueue)Size() int  {
	return len(oq.nodes)
}


func (oq *OrderQueue) GeAllNodes() []interface{} {
	oq.lock.Lock()
	defer oq.lock.Unlock()

	var nodes []interface{}

	for i := 0; i < len(oq.nodes); i++ {
		nodes = append(nodes, oq.nodes[i])
	}

	return nodes
}
