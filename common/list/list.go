package list

import (
	"fmt"
	"github.com/kprc/nbsdht/nbslink"
	"sync/atomic"
	"github.com/pkg/errors"
)

type list struct {
	root *nbslink.LinkNode
	cnt int32
	cmp func(v1 interface{},v2 interface{}) int   //return 0 is equal, for get node
	clone func(v1 interface{}) (r interface{})		//for copy
	sort func(v1 interface{},v2 interface{}) int    //for sort
}

type FDel func(arg interface{}, v interface{}) bool  // if true, delete the node
type FDo func(arg interface{}, v interface{}) (ret interface{},err error)

type ListCusor struct {
	arrv []interface{}
	cursor int
}

type List interface {
	Add(node *nbslink.LinkNode)
	Append(node *nbslink.LinkNode)
	Del(node *nbslink.LinkNode)
	UpdateValueOrder(v interface{})
	Find(v interface{}) *nbslink.LinkNode
	FindDo(arg interface{},do FDo) (ret interface{},err error)
	AddValue(v interface{})
	AppendValue(v interface{})
	DelValue(v interface{})
	Count() int32
	Traverse(arg interface{},fDo FDo)
	TraverseDel(arg interface{},fDel FDel)
	Clone() List
	DeepClone() (List,error)
	ConCat(l List,sort bool)
	DeepConCat(l List, sort bool) error
	SetCloneFunc(clone func(v1 interface{}) (r interface{}))
	SetSortFunc(sort func(v1 interface{},v2 interface{}) int )
	AddValueOrder(v interface{})
	GetClone() (clone func(v1 interface{}) (r interface{}))
	GetSortFunc() (sort func(v1 interface{},v2 interface{}) int )
	ListIterator(cnt int) *ListCusor
	Sort()
	GetFirst() (v interface{},err error)
	GetLast()  (v interface{},err error)
	InsertBefore(v interface{},before interface{}) error
	InsertAfter(v interface{},after interface{}) error
}

func NewList(cmp func(v1 interface{},v2 interface{}) int) List  {
	if cmp == nil {
		return nil
	}
	return &list{cmp:cmp}
}

func Print(arg interface{},data interface{}) (ret interface{},err error) {
	s:=arg.(string)
	idata:=data.(int)

	fmt.Println(s,idata)


	return
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

func (l *list)GetFirst() (v interface{},err error)  {

	if l.root == nil{
		return nil,errors.New("list is empty")
	}

	return l.root.Value,nil
}

func (l *list)GetLast() (v interface{},err error){
	if l.root == nil{
		return nil,errors.New("list is empty")
	}
	return l.root.Prev().Value,nil
}

func (l *list)InsertBefore(v interface{},before interface{}) error {
	if v==nil || before == nil{
		return errors.New("parameter error")
	}
	n:=l.Find(before)
	if n == nil{
		return errors.New("before not found")
	}
	nv:= nbslink.NewLink(v)
	n.Insert(nv)
	if n == l.root{
		l.root = nv
	}
	l.incCnt()

	return nil
}

func (l *list)InsertAfter(v interface{},after interface{}) error {
	if v==nil || after == nil{
		return errors.New("parameter error")
	}
	n:=l.Find(after)
	if n == nil{
		return errors.New("after not found")
	}
	nv:= nbslink.NewLink(v)
	n.Add(nv)

	l.incCnt()

	return nil
}

func (l *list)Append(node *nbslink.LinkNode)  {
	if node == nil{
		return
	}
	if l.root == nil{
		node.Init()
		l.root = node
	}else{
		l.root.Add(node)
		l.incCnt()
	}
}

func (l *list)AppendValue(v interface{})  {
	if v == nil{
		return
	}
	n:= nbslink.NewLink(v)
	l.Append(n)
}

func (l *list)ListIterator(cnt int) *ListCusor  {
	lc:=&ListCusor{}
	lc.arrv = make([]interface{},0)

	if l.root == nil{
		return lc
	}
	root:=l.root
	node := l.root

	curcnt:=0

	for  {
		if cnt <= 0 || curcnt < cnt{
			v:=node.Value
			if l.clone != nil{
				v=l.clone(node.Value)
			}
			lc.arrv = append(lc.arrv,v)
			curcnt ++
		}else{
			break
		}
		node = node.Next()
		if node == root {
			break
		}
	}

	return lc

}

func (it *ListCusor)Next() interface{} {
	if it.cursor >= len(it.arrv){
		return nil
	}

	v:=it.arrv[it.cursor]
	it.cursor ++

	return v
}

func (it *ListCusor)Reset() {
	it.cursor = 0
}

func (it *ListCusor)Count() int {
	return len(it.arrv)
}

func (l *list)Sort()  {
	lnew:=NewList(l.cmp)
	lnew.SetCloneFunc(l.clone)
	lnew.SetSortFunc(l.sort)

	lnew.ConCat(l,true)

	*l = *(lnew.(*list))
}



func (l *list)incCnt()  {
	atomic.AddInt32(&l.cnt,1)
}

func (l *list)decCnt()  {
	atomic.AddInt32(&l.cnt,-1)
}

func (l *list)GetClone() (clone func(v1 interface{}) (r interface{})) {
	return l.clone
}

func (l *list)GetSortFunc() (sort func(v1 interface{},v2 interface{}) int )  {
	return l.sort
}

func (l *list)SetCloneFunc(clone func(v1 interface{}) (r interface{}))  {
	l.clone = clone
}

func (l *list)SetSortFunc(sort func(v1 interface{},v2 interface{}) int ) {
	l.sort = sort
}

func (l *list)ConCat(nl List,sort bool)  {
	if nl==nil || nl.Count() == 0{
		return
	}

	nnl:=nl.(*list)

	root:=nnl.root
	node:=nnl.root

	for{
		if sort{
			l.AddValueOrder(node.Value)
		}else{
			l.AddValue(node.Value)
		}
		node=node.Next()
		if node == root{
			break
		}
	}
}

func (l *list)UpdateValueOrder(v interface{})  {
	l.DelValue(v)
	l.AddValueOrder(v)
}


func (l *list)DeepConCat(nl List, sort bool) error  {
	if nl==nil || nl.Count() == 0 {
		return nil
	}
	if nl.GetClone() == nil{
		return errors.New("No clone function")
	}

	nnl:=nl.(*list)

	root:=nnl.root
	node:=nnl.root

	for{
		if sort {
			l.AddValueOrder(nnl.clone(node.Value))
		}else{
			l.AddValue(nnl.clone(node.Value))
		}
		node=node.Next()
		if node == root{
			break
		}
	}
	return nil
}

func (l *list)Clone() List  {
	newl:=NewList(l.cmp)

	if l.root == nil{
		return newl
	}

	root:=l.root
	node := l.root

	for  {
		newl.AddValue(node.Value)
		node = node.Next()
		if node == root {
			break
		}
	}

	return newl
}

func (l *list)DeepClone() (List,error)  {
	if l.clone == nil{
		return nil,errors.New("need clone function, please set it")
	}

	newl:=NewList(l.cmp)
	if l.root == nil{
		return newl,nil
	}

	root:=l.root
	node := l.root

	for  {
		newl.AddValue(l.clone(node.Value))
		node = node.Next()
		if node == root {
			break
		}
	}

	return newl,nil
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
	if l.root == nil{
		return
	}
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

func (l *list)AddValueOrder(v interface{})  {
	node:= nbslink.NewLink(v)

	if l.root == nil {
		node.Init()
		l.root = node
		l.incCnt()
		return
	}
	nxt:=l.root
	prev:=nxt

	for{
		d:=l.sort(v,nxt.Value)
		if d==0{
			break
		}
		if d<0{
			nxt.Insert(node)
			l.incCnt()
			if nxt==l.root{
				l.root = node
			}
			break
		}
		prev=nxt
		nxt=nxt.Next()

		if nxt == l.root{
			prev.Add(node)
			l.incCnt()
			break
		}
	}
}

func (l *list)DelValue(v interface{})  {

	if l.root == nil{
		return
	}
	root:=l.root
	n:=l.root

	for {
		nxt:=n.Next()
		if 0 == l.cmp(v,n.Value) {
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
		n = nxt
		if n == root {
			break
		}

	}
}

func (l *list)Count() int32 {
	return atomic.LoadInt32(&l.cnt)
}

func (l *list)Traverse(arg interface{},fDo FDo)  {
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

func (l *list)Find(v interface{}) *nbslink.LinkNode  {
	if l.root == nil{
		return nil
	}

	root:=l.root
	n:=l.root


	for {
		if 0 == l.cmp(v,n.Value){
			return n
		}
		n = n.Next()
		if n == root {
			break
		}

	}
	return nil
}

func (l *list)FindDo(arg interface{},do FDo) (ret interface{},err error)  {
	if l.root == nil{
		return nil,errors.New("list is none")
	}

	root:=l.root
	n:=l.root


	for {
		if 0 == l.cmp(arg,n.Value){
			return do(arg,n.Value)
		}
		n = n.Next()
		if n == root {
			break
		}

	}
	return nil,errors.New("not found node")
}

func (l *list)TraverseDel(arg interface{},fDel FDel)  {
	root:=l.root
	n:=l.root

	for {
		nxt:=n.Next()
		if fDel(arg,n.Value) {

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
		}
		n = nxt
		if n == root {
			break
		}

	}
}