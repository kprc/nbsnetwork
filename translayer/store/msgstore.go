package store

import (
	"github.com/kprc/nbsnetwork/common/list"
	"sync"
)

type block struct {
	blk interface{}
	timeout int
	sendtime int64
	resndCnt int
	sn uint64
}






