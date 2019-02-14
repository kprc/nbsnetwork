package queue

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsdht/nbslink"
)

var(
	LEVEL0 = 0
	LEVEL1 = 1
	LEVEL2 = 2
	LEVEL3 = 3
	LEVEL4 = 4
	LEVEL5 = 5
)


type priorityQueue struct {
	pq [6]Queue
}


type PriorityQueue interface {
	EnQueue(level int,node *nbslink.LinkNode) error
	DeQueue() *nbslink.LinkNode
	EnQueueValue(level int,v interface{}) error
	DeQueueValue() interface{}
	Traverse(arg interface{},fDo func(arg interface{},data interface{}))
}

func NewPriorityQueue() PriorityQueue  {
	pq := &priorityQueue{}

	return pq
}

func (pq *priorityQueue) EnQueue(level int,node *nbslink.LinkNode) error {
	if level < 0 || level >=len(pq.pq) {
		return nbserr.NbsErr{ErrId:nbserr.ERROR_DEFAULT,Errmsg:"level error"}
	}

	if pq.pq[level] == nil{
		pq.pq[level] = NewQueue()
	}

	q:=pq.pq[level]


	q.EnQueue(node)

	return nil
}

func (pq *priorityQueue) DeQueue() *nbslink.LinkNode {
	for i:=0; i<len(pq.pq) ;i++  {
		q:=pq.pq[i]
		node := q.DeQueue()
		if node != nil {
			return node
		}
	}

	return nil
}

func (pq *priorityQueue) EnQueueValue(level int,v interface{}) error {
	if level < 0 || level >=len(pq.pq) {
		return nbserr.NbsErr{ErrId:nbserr.ERROR_DEFAULT,Errmsg:"level error"}
	}

	if pq.pq[level] == nil{
		pq.pq[level] = NewQueue()
	}

	q:=pq.pq[level]

	q.EnQueueValue(v)

	return nil

}

func (pq *priorityQueue) DeQueueValue() interface{} {
	node := pq.DeQueue()

	if node == nil {
		return nil
	}

	return node.Value
}

func (pq *priorityQueue) Traverse(arg interface{},fDo func(arg interface{},data interface{})) {
	for i:=0;i<len(pq.pq); i++{
		q:=pq.pq[i]
		q.Traverse(arg,fDo)
	}
}



