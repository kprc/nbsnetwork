package applayer

import (
	"github.com/kprc/nbsdht/nbserr"
	"sync"
)

type Appdo func(rcv interface{}, arg interface{}) (v interface{}, err error)

type appblock struct {
	do Appdo
}

type appblockstore struct {
	handle map[uint32]*appblock
}

type AppBlockStore interface {
	Reg(apptyp uint32, do Appdo)
	Do(apptyp uint32, rcv interface{}, arg interface{}) (v interface{}, err error)
}

var (
	appblockstoreinstance AppBlockStore
	pbsinstancelock       sync.Mutex
	appparamerr           = nbserr.NbsErr{ErrId: nbserr.APP_TYP_ERR, Errmsg: "Not found apptyp"}
)

func GetAppBlockStore() AppBlockStore {
	if appblockstoreinstance != nil {
		return appblockstoreinstance
	}
	pbsinstancelock.Lock()
	defer pbsinstancelock.Unlock()
	if appblockstoreinstance == nil {
		appblockstoreinstance = newAppBlockStore()
	}

	return appblockstoreinstance

}

func newAppBlockStore() AppBlockStore {
	return &appblockstore{handle: make(map[uint32]*appblock)}
}

func (aps *appblockstore) Reg(apptyp uint32, do Appdo) {
	ab := &appblock{do: do}

	aps.handle[apptyp] = ab
}

func (aps *appblockstore) Do(apptyp uint32, rcv interface{}, arg interface{}) (v interface{}, err error) {
	if pb, ok := aps.handle[apptyp]; !ok {
		return nil, appparamerr
	} else {
		return pb.do(rcv, arg)
	}
}
