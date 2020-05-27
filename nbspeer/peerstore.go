package nbspeer

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"sync"
)

type peerstore struct {
	hashlist.HashList
	tick chan int64
}

type PeerStore interface {
}

var (
	psinstancelock sync.Mutex
	psinstance     PeerStore
	lasttimeout    int64
	timeouttv      int64 = 1000 //ms
)

var fhash = func(v interface{}) {
	p := v.(Peer)
	uid := p.GetUid()
}
