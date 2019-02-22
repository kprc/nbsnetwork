package list

import (
	"fmt"
	"github.com/kprc/nbsdht/nbslink"
	"sync/atomic"
)

type list struct {
	root *nbslink.LinkNode
	cnt int32
	cmp func(v1 interface{},v2 interface{}) int   //return 0 is equal
}

type List interface {
	Add(node *nbslink.LinkNode)
	Del(node *nbslink.LinkNode)
	AddValue(v interface{})
	DelValue(v interface{})
	Count() int32
	Traverse(arg interface{},fDo func(arg interface{},data interface{}))
}

func NewList(cmp func(v1 interface{},v2 interface{}) int) List  {
	if cmp == nil {
		return nil
	}
	return &list{cmp:cmp}
}

func Print(arg interface{},data interface{})  {
	s:=arg.(string)
	idata:=data.(int)

	fmt.Println(s,idata)
}

func Cmp(v1 interface{},v2 interface{})  int {
	i1 := v1.(int)
	i2 := v2.(int)

	if i1==i2 {
		return 0
	}

	if i1 > i2 {
		return 1
	}else {
		return -1
	}
}


func (l *list)incCnt()  {
	atomic.AddInt32(&l.cnt,1)
}

func (l *list)decCnt()  {
	atomic.AddInt32(&l.cnt,-1)
}

func (l *list)Add(node *nbslink.LinkNode)  {
	if node == nil {
		return
	}
	if l.root == nil {
		node.Init()
		l.root = node
		l.incCnt()
		return
	}

	l.incCnt()
	l.root.Insert(node)
}

func (l *list)Del(node *nbslink.LinkNode)  {
	root:=l.root
	n:=l.root

	for {
		if n==node {
			nxt := n.Next()
			if n == root{
				l.root =nxt
				if nxt == n{
					l.root = nil
					l.decCnt()
					return
				}
			}

			n.Remove()
			l.decCnt()
			break
		}
		n = n.Next()
		if n == root {
			break
		}
	}

}

func (l *list)AddValue(v interface{})  {
	n:= nbslink.NewLink(v)
	l.Add(n)
}

func (l *list)DelValue(v interface{})  {
	root:=l.root
	n:=l.root

	for {
		if 0 == l.cmp(v,n.Value) {
			nxt := n.Next()
			if n == root{
				l.root =nxt
				if nxt == n{
					l.root = nil
					l.decCnt()
					return
				}
			}

			n.Remove()
			l.decCnt()
			break
		}
		n = n.Next()
		if n == root {
			break
		}

	}
}

func (l *list)Count() int32 {
	return atomic.LoadInt32(&l.cnt)
}

func (l *list)Traverse(arg interface{},fDo func(arg interface{},data interface{}))  {
	if l.root == nil{
		return
	}

	root:=l.root
	node := l.root

	for  {
		fDo(arg,node.Value)
		node = node.Next()
		if node == root {
			break
		}
	}
}