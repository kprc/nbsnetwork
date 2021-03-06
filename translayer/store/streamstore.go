package store

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type streamblk struct {
	key             UdpStreamKey
	blk             interface{}
	timeoutInterval int32
	lastAccessTime  int64
}

type streamstore struct {
	hashlist.HashList
	tick chan int64
	quit chan int64
	wg   *sync.WaitGroup
}

type StreamStore interface {
	AddStream(s interface{})
	AddStreamWithParam(s interface{}, timeInterval int32)
	DelStream(s interface{})
	FindStreamDo(s interface{}, arg interface{}, do list.FDo) (r interface{}, err error)
	Run()
	Stop()
}

var (
	ssInstanceLock   sync.Mutex
	ssIntance        StreamStore
	ssLastAccessTime int64
	ssTimeOutTV      int64 = 1000 //ms
)

func (sb *streamblk) GetKey() UdpStreamKey {
	return sb.key
}

func GetStreamStoreInstance() StreamStore {
	if ssIntance == nil {
		ssInstanceLock.Lock()
		defer ssInstanceLock.Unlock()

		if ssIntance == nil {
			ssIntance = newStreamStore()
		}
	}

	return ssIntance
}

func newStreamStore() StreamStore {

	hl := hashlist.NewHashList(0x80, StreamKeyHash, StreamKeyEquals)
	ss := &streamstore{hl, make(chan int64, 64), make(chan int64, 1), &sync.WaitGroup{}}

	t := tools.GetNbsTickerInstance()
	t.Reg(&ss.tick)

	ssLastAccessTime = tools.GetNowMsTime()

	return ss
}

func GetStreamBlkAndRefresh(v interface{}) interface{} {
	sb := v.(*streamblk)
	sb.lastAccessTime = tools.GetNowMsTime()
	return sb.blk
}

func GetStreamBlk(v interface{}) interface{} {
	sb := v.(*streamblk)
	return sb.blk
}

func (ss *streamstore) addStream(v interface{}, timeinterval int32) {
	sb := &streamblk{blk: v}
	sb.key = v.(StreamKeyInter).GetKey()
	sb.timeoutInterval = timeinterval
	sb.lastAccessTime = tools.GetNowMsTime()

	ss.Add(sb)
}

func (ss *streamstore) AddStream(v interface{}) {
	ss.addStream(v, int32(constant.UDP_STREAM_STORE_TIMEOUT))
}

func (ss *streamstore) AddStreamWithParam(v interface{}, timeInterval int32) {
	if timeInterval == 0 {
		timeInterval = int32(constant.UDP_STREAM_STORE_TIMEOUT)
	}
	ss.addStream(v, timeInterval)
}

func (ss *streamstore) DelStream(v interface{}) {
	ss.Del(v)
}

func (ss *streamstore) FindStreamDo(s interface{}, arg interface{}, do list.FDo) (r interface{}, err error) {
	return ss.FindDo(s, arg, do)
}

func (ss *streamstore) doTimeOut() {
	type s2del struct {
		arrdel []*streamblk
	}

	arr2del := &s2del{arrdel: make([]*streamblk, 0)}

	fdo := func(arg interface{}, v interface{}) (r interface{}, err error) {
		sb := v.(*streamblk)
		l := arg.(*s2del)

		curtime := tools.GetNowMsTime()
		tv := curtime - sb.lastAccessTime
		if tv > int64(sb.timeoutInterval) {
			l.arrdel = append(l.arrdel, sb)
		}
		return
	}

	ss.TraversAll(arr2del, fdo)

	for _, sb := range arr2del.arrdel {
		ss.DelStream(sb)
	}

}

func (ss *streamstore) Run() {
	ss.wg.Add(1)
	defer ss.wg.Done()
	fmt.Println("Stream Store is Running")
	for {
		select {
		case <-ss.tick:
			if tools.GetNowMsTime()-ssLastAccessTime > ssTimeOutTV {
				ssLastAccessTime = tools.GetNowMsTime()
				ss.doTimeOut()
			}
		case <-ss.quit:
			return
		}
	}
}

func (ss *streamstore) Stop() {
	ss.quit <- 1
	if ss.wg != nil {
		ss.wg.Wait()
	}
}
