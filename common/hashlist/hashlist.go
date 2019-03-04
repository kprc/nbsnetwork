package hashlist

import (
	"github.com/kprc/nbsnetwork/common/list"
	"sync"
)

type hashlist struct {
	hashcnt uint
	bucket []list.List
	bucketlock []*sync.Mutex
	fhash func(v interface{}) uint
	fequals func(v1 interface{},v2 interface{}) int
	gcnt int64
}

type HashList interface {
	Add(v interface{})
	Del(v interface{})
	FindDo(v interface{},arg interface{}, do FDo) (ret interface{},err error)
}

type FDo func(arg interface{}, v interface{}) (ret interface{},err error)

func NewHashList(hashcnt uint,fhash func(key interface{}) uint,fequals func(v1 interface{},v2 interface{}) int) HashList  {
	hl := &hashlist{}



	for i,_:=range hl.bucket{
		hl.bucket[i] = list.NewList(fequals)
	}
	for i,_:=range hl.bucketlock{
		hl.bucketlock[i] = &sync.Mutex{}
	}

	hl.fequals = fequals
	hl.fhash =fhash

	return hl

}

func (hl *hashlist)Add(v interface{})  {
	hash:=hl.fhash(v)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()


	hl.bucket[hash].AddValue(v)
	hl.gcnt ++

}

func (hl *hashlist)Del(v interface{})  {
	hash:=hl.fhash(v)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()


	hl.bucket[hash].DelValue(v)
	cnt :=hl.bucket[hash].Count()
	cnt = cnt - hl.bucket[hash].Count()

	hl.gcnt = hl.gcnt - int64(cnt)
}

func (hl *hashlist)FindDo(v interface{},arg interface{},do  FDo ) (ret interface{},err error)  {
	hash:=hl.fhash(v)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()

	root:=hl.bucket[hash]
	node := root.Find(v)

	if node == nil {
		return
	}

	return do(arg,node.Value)
}
