package queue

import (
	"fmt"
	"github.com/kprc/nbsdht/nbslink"
)

type queue struct {
	root *nbslink.LinkNode
}



type Queue interface {
	EnQueue(node *nbslink.LinkNode)
	DeQueue() *nbslink.LinkNode
	EnQueueValue(v interface{})
	DeQueueValue() interface{}
	Traverse(arg interface{},fDo func(arg interface{},data interface{}))
}

func NewQueue() Queue {
	return &queue{}
}

func Print(arg interface{},data interface{})  {
	s:=arg.(string)
	sdata := data.(string)

	fmt.Println(s,sdata)
}

func (q *queue) EnQueue(node *nbslink.LinkNode) {
	if node == nil{
		return
	}

	if q.root == nil {
		node.Init()
		q.root = node
		return
	}

	q.root.Insert(node)
}

func (q *queue)DeQueue() *nbslink.LinkNode  {
	if q.root == nil {
		return nil
	}

	node:=q.root

	q.root = q.root.Next()

	if q.root == node {
		q.root = nil
	}else {
		node.Remove()
	}
	return node
}

func (q *queue)EnQueueValue(v interface{}){
	node := nbslink.NewLink(v)
	q.EnQueue(node)
}
func (q *queue)DeQueueValue() interface{}{
	node := q.DeQueue()
	if node == nil {
		return nil
	}
	return node.Value
}

func (q *queue)Traverse(arg interface{},fDo func(arg interface{},data interface{}))  {
	if q.root == nil{
		return
	}

	root:=q.root
	node := q.root

	for  {
		fDo(arg,node.Value)
		node = node.Next()
		if node == root {
			break
		}
	}
}
