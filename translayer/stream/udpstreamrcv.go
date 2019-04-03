package stream

import (
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
	"sync"
)

type streamrcv struct {
	udpmsgcache map[uint64]store.UdpMsg
	lastreadpos uint64
	finishflag bool
	toppos uint64
	sn uint64
	uid string
	w io.Writer
}

type StreamRcv interface {
	Recv()
}

type sskey struct {
	uid string
	sn uint64
}

type streamstore struct {
	lock sync.Mutex
	store map[sskey]StreamRcv
	storelock map[sskey]sync.Mutex
}


type StreamStore interface {

}
