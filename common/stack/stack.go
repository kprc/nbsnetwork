package stack

import (
	"github.com/kprc/nbsdht/nbslink"
	"sync/atomic"
	"fmt"
)

type stack struct {
	root *nbslink.LinkNode
	cnt int32
}

type Stack interface {
	Push(v interface{})
	Pop() interface{}
	Count() int32
	Traverse(arg interface{},fDo func(arg interface{},data interface{}))
}

func NewStack() Stack  {
	return &stack{root:nbslink.NewLink(nil).Init()}
}

func Print(arg interface{},data interface{})  {
	s:=arg.(string)
	sdata := data.(string)

	fmt.Println(s,sdata)
}

func (s *stack)Push(v interface{}) {
	l:=nbslink.NewLink(v)

	s.root.Add(l)
	atomic.AddInt32(&s.cnt,1)
}

func (s *stack)Pop() interface{}  {
	nxt := s.root.Next()

	if nxt == s.root {
		return nil
	}
	nxt.Remove()
	atomic.AddInt32(&s.cnt,-1)
	return nxt.Value
}


func (s *stack)Count() int32{
	return atomic.LoadInt32(&s.cnt)
}

func (s *stack)Traverse(arg interface{},fDo func(arg interface{},data interface{}))  {
	if s.root == nil || s.root.Next() == s.root{
		return
	}

	root := s.root
	node := s.root.Next()

	for  {
		if node == root {
			break
		}
		fDo(arg,node.Value)
		node = node.Next()
	}
}







