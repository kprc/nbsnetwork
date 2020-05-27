package hdb

import (
	"github.com/kprc/nbsnetwork/common/list"
)

type FixedQueueIntf interface {
	EnQueue(v interface{})
	Iterator() *list.ListCusor
	QSize() int
	CurSize() int
	GetTopN(begin, topn int) []interface{}
}

type FixedQueue struct {
	qSize   int
	curSize int
	qV      list.List
}

func NewFixedQueue(size int, cmp func(v1 interface{}, v2 interface{}) int) FixedQueueIntf {
	fq := &FixedQueue{}
	if size <= 0 {
		size = 100
	}
	fq.qSize = size

	fq.qV = list.NewList(cmp)

	return fq
}

func (fq *FixedQueue) EnQueue(v interface{}) {
	fq.qV.InsertValue(v)
	fq.curSize++
	for {
		if fq.curSize > fq.qSize {
			node, _ := fq.qV.GetLastNode()
			fq.qV.Del(node)
			fq.curSize--
		} else {
			break
		}
	}
}

func (fq *FixedQueue) Iterator() *list.ListCusor {
	return fq.qV.ListIterator(fq.curSize)
}

func (fq *FixedQueue) QSize() int {
	return fq.qSize
}

func (fq *FixedQueue) CurSize() int {
	return fq.curSize
}

func (fq *FixedQueue) GetTopN(begin, topn int) []interface{} {

	type arrinterface struct {
		arr         []interface{}
		begin, topn int
	}

	var arg *arrinterface = &arrinterface{}

	fq.qV.Traverse(arg, func(arg interface{}, v interface{}) (ret interface{}, err error) {
		parr := arg.(*arrinterface)

		if v.(*HDBV).Cnt >= parr.begin && v.(*HDBV).Cnt < parr.begin+parr.topn {
			parr.arr = append(parr.arr, v)
		}

		return

	})

	return arg.arr

}
