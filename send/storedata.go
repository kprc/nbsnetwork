package send

import (
	"sync/atomic"
)

type storeData struct {
	bdr BlockDataer
	cnt int32
}


type StoreDataer interface {
	ReferCntInc()
	ReferCntDec()
	GetReferCnt() int32
}

func NewStoreData(bdr BlockDataer)StoreDataer {
	return &storeData{bdr:bdr}
}

func (sd *storeData)ReferCntInc()  {
	atomic.AddInt32(&sd.cnt,1)
}

func (sd *storeData)ReferCntDec()  {
	atomic.AddInt32(&sd.cnt,-1)
}

func (sd *storeData)GetReferCnt() int32  {
	return atomic.LoadInt32(&sd.cnt)
}

