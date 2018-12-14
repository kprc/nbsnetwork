package nbsnetwork

import (
	"sync/atomic"
)

type storeData struct {
	bdr BlockDataer
	cnt uint32
}


type StoreDataer interface {
	ReferCntInc()
	ReferCntDec()
	GetReferCnt() uint32
}

func NewStoreData(bdr BlockDataer)StoreDataer {
	return &storeData{bdr:bdr}
}

func (sd *storeData)ReferCntInc()  {
	atomic.AddUint32(&sd.cnt,1)
}

func (sd *storeData)ReferCntDec()  {
	atomic.AddUint32(&sd.cnt,-1)
}

func (sd *storeData)GetReferCnt() uint32  {
	return sd.cnt
}

