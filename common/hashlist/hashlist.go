package hashlist

import (
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type FHash func(v interface{}) uint
type Fequals func(v1 interface{},v2 interface{}) int


type hashlist struct {
	bucketsize uint
	realbucketsize uint
	bucket []list.List
	bucketlock []*sync.Mutex
	fhash FHash
	fequals Fequals
	gcnt int64
}

type HashList interface {
	Add(v interface{})
	Del(v interface{})
	FindDo(v interface{},arg interface{}, do FDo) (ret interface{},err error)
}

type FDo func(arg interface{}, v interface{}) (ret interface{},err error)

func NewHashList(bucketsize uint,fhash FHash,fequals Fequals) HashList  {
	hl := &hashlist{bucketsize:bucketsize,fhash:fhash,fequals:fequals}

	rbz:= tools.ResizeHash(bucketsize)
	hl.realbucketsize = rbz

	listbucket := make([]list.List,rbz)

	for i,_:=range listbucket{
		listbucket[i] = list.NewList(fequals)
	}

	hl.bucket = listbucket

	lockbucket:=make([]*sync.Mutex,rbz)

	for i,_:=range lockbucket{
		lockbucket[i] = &sync.Mutex{}
	}

	hl.bucketlock = lockbucket

	return hl

}

func (hl *hashlist)Add(v interface{})  {
	hash:=hl.fhash(v)
	hash = hash & (hl.realbucketsize - 1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()


	hl.bucket[hash].AddValue(v)
	hl.gcnt ++

}

func (hl *hashlist)Del(v interface{})  {
	hash:=hl.fhash(v)
	hash = hash & (hl.realbucketsize - 1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()

	cnt :=hl.bucket[hash].Count()
	hl.bucket[hash].DelValue(v)

	cnt = cnt - hl.bucket[hash].Count()

	hl.gcnt = hl.gcnt - int64(cnt)
}

func (hl *hashlist)FindDo(v interface{},arg interface{},do  FDo ) (ret interface{},err error)  {
	hash:=hl.fhash(v)
	hash = hash & (hl.realbucketsize - 1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()

	root:=hl.bucket[hash]
	node := root.Find(v)

	if node == nil {
		return
	}

	return do(arg,node.Value)
}
