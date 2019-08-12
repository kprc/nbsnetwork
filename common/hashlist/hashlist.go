package hashlist

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"math/big"
)

type FHash func(v interface{}) uint
type Fequals func(v1 interface{},v2 interface{}) int
type FGetDistance func(v interface{}) big.Int


type hashlist struct {
	bucketsize uint
	realbucketsize uint
	bucket []list.List
	bucketlock []*sync.Mutex
	fhash FHash
	fequals Fequals
	gcnt int64
	limitlistcnt int32
	clone func(v1 interface{}) (r interface{})
	sort func(v1 interface{},v2 interface{}) int
}

type HashList interface {
	Add(v interface{})
	AddOrder(v interface{})
	UpdateOrder(v interface{})
	Del(v interface{})
	GetNodeCntOneBucket(v interface{}) int32
	GetNodes(v interface{},nCnt int32,sort func(v1 interface{},v2 interface{}) int) list.List
	FindDo(v interface{},arg interface{}, do list.FDo) (ret interface{},err error)
	TraversAll(arg interface{}, do list.FDo)
	SetCloneFunc(clone func(v1 interface{}) (r interface{}))
	SetSortFunc(sort func(v1 interface{},v2 interface{}) int)
	SetLimitCnt(n int32)
	GetLimitCnt() int32
}



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

func (hl *hashlist)SetLimitCnt(n int32) {
	hl.limitlistcnt = n
}

func (hl *hashlist)GetLimitCnt()  int32{
	return hl.limitlistcnt
}

func (hl *hashlist)SetCloneFunc(clone func(v1 interface{}) (r interface{})) {
	hl.clone = clone
	var i int
	for ;i<int(hl.realbucketsize);i++{
		hl.bucket[i].SetCloneFunc(clone)
	}
}

func (hl *hashlist)SetSortFunc(sort func(v1 interface{},v2 interface{}) int)  {
	hl.sort = sort
	var i int
	for ;i<int(hl.realbucketsize);i++{
		hl.bucket[i].SetSortFunc(sort)
	}
}

func (hl *hashlist)Add(v interface{})  {
	hash:=hl.fhash(v)
	hash = hash & (hl.realbucketsize - 1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()
	if hl.limitlistcnt > 0 && hl.bucket[hash].Count() > hl.limitlistcnt{
		return
	}

	hl.bucket[hash].AddValue(v)
	hl.gcnt ++

}

func (hl *hashlist)AddOrder(v interface{})  {
	hash:=hl.fhash(v)
	hash = hash & (hl.realbucketsize - 1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()

	if hl.limitlistcnt > 0 && hl.bucket[hash].Count() > hl.limitlistcnt{
		return
	}

	hl.bucket[hash].AddValueOrder(v)
	hl.gcnt ++

}

func (hl *hashlist)UpdateOrder(v interface{})  {
	hash:=hl.fhash(v)
	hash = hash & (hl.realbucketsize - 1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()

	if v:=hl.bucket[hash].Find(v);v!=nil{
		hl.bucket[hash].UpdateValueOrder(v)
		return
	}
	if hl.limitlistcnt > 0 && hl.bucket[hash].Count() > hl.limitlistcnt{
		return
	}
	hl.bucket[hash].AddValueOrder(v)
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

func (hl *hashlist)GetNodeCntOneBucket(v interface{}) int32  {
	hash:=hl.fhash(v)

	hash = hash & (hl.realbucketsize -1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()

	return hl.bucket[hash].Count()
}

func (hl *hashlist)GetNodes(v interface{},nCnt int32,sort func(v1 interface{},v2 interface{}) int) list.List  {
	hash:=hl.fhash(v)

	hash = hash & (hl.realbucketsize -1)

	l:=list.NewList(hl.fequals)

	l.SetSortFunc(sort)
	l.SetCloneFunc(hl.clone)

	nnCnt:=int32(hl.limitlistcnt)
	if nCnt > 0{
		nnCnt = nCnt
	}

	i:=int(hash)
	for ;i>=0;i--{
		hl.bucketlock[i].Lock()
		l.DeepConCat(hl.bucket[i],true)
		hl.bucketlock[i].Unlock()
		if l.Count() >= nnCnt{
			break
		}
	}
	i=int(hash)+1
	if l.Count()<nnCnt && i<int(hl.realbucketsize){
		for ; i<int(hl.realbucketsize); i++{
			hl.bucketlock[i].Lock()
			l.DeepConCat(hl.bucket[i],true)
			hl.bucketlock[i].Unlock()
			if l.Count() >= nnCnt{
				break
			}
		}
	}

	return l
}


func (hl *hashlist)FindDo(v interface{},arg interface{},do  list.FDo ) (ret interface{},err error)  {
	hash:=hl.fhash(v)
	hash = hash & (hl.realbucketsize - 1)

	hl.bucketlock[hash].Lock()
	defer hl.bucketlock[hash].Unlock()

	root:=hl.bucket[hash]
	node := root.Find(v)

	if node == nil {
		return nil,nbserr.NbsErr{ErrId:nbserr.HASHLIST_NO_NODE_ERR,Errmsg:"Node not found"}
	}

	return do(arg,node.Value)
}


func (hl *hashlist)TraversAll(arg interface{}, do list.FDo)  {

	var ui uint
	for ;ui<hl.realbucketsize ;ui++  {
		hl.bucketlock[ui].Lock()
		hl.bucket[ui].Traverse(arg,do)
		hl.bucketlock[ui].Unlock()
	}

}
