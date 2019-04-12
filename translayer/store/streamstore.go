package store

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"sync"
	"github.com/kprc/nbsnetwork/tools"
)


type streamblk struct {
	key UdpStreamKey
	blk interface{}
	timeoutInterval int32
	lastAccessTime int64
}




type streamstore struct {
	hashlist.HashList
	tick chan int64
}

type StreamStore interface {
	AddStream(s interface{})
	AddStreamWithParam(s interface{}, timeInterval int32)
	DelStream(s interface{})
	FindStreamDo(s interface{},arg interface{}, do list.FDo) (r interface{}, err error)
	Run()
}

var (
	ssInstanceLock sync.Mutex
	ssIntance StreamStore
	ssLastAccessTime int64
	ssTimeOutTV int64 = 1000   //ms
)

func (sb *streamblk)GetKey() UdpStreamKey {
	return sb.key
}

func GetStreamStoreInstance() StreamStore  {
	if ssIntance == nil{
		ssInstanceLock.Lock()
		defer ssInstanceLock.Unlock()

		if ssIntance == nil{
			ssIntance = newStreamStore()
		}
	}

	return ssIntance
}

func newStreamStore() StreamStore  {

	hl:=hashlist.NewHashList(0x80,StreamKeyHash,StreamKeyEquals)
	ss:=&streamstore{hl,make(chan int64,64)}

	t:=tools.GetNbsTickerInstance()
	t.Reg(&ss.tick)

	ssLastAccessTime = tools.GetNowMsTime()

	return ss
}

func GetStreamBlkAndRefresh(v interface{}) interface{}{
	sb:=v.(streamblk)
	sb.lastAccessTime = tools.GetNowMsTime()
	return sb.blk
}

func GetStreamBlk(v interface{}) interface{}  {
	sb:=v.(streamblk)
	return sb.blk
}


func (ss *streamstore)addStream(v interface{}, timeinterval int32)  {
	sb:=&streamblk{blk:v}
	sb.key = v.(StreamKeyInter).GetKey()
	sb.timeoutInterval = timeinterval
	sb.lastAccessTime = tools.GetNowMsTime()

	ss.Add(sb)
}

func (ss *streamstore)AddStream(v interface{})  {
	ss.addStream(v,10000)
}

func (ss *streamstore)AddStreamWithParam(v interface{},timeInterval int32)  {
	ss.addStream(v,timeInterval)
}

func (ss *streamstore)DelStream(v interface{})  {
	ss.Del(v)
}

func (ss *streamstore)FindStreamDo(s interface{},arg interface{}, do list.FDo) (r interface{}, err error)  {
	return ss.FindDo(s,arg,do)
}

func (ss *streamstore)doTimeOut()  {
	type s2del struct {
		arrdel []*streamblk
	}

	arr2del := &s2del{arrdel:make([]*streamblk,0)}

	fdo:=func(arg interface{}, v interface{})(r interface{},err error){
		sb:=v.(*streamblk)
		l:=arg.(*s2del)

		curtime := tools.GetNowMsTime()
		tv := curtime - sb.lastAccessTime
		if tv > int64(sb.timeoutInterval) {
			l.arrdel = append(l.arrdel,sb)
		}
		return
	}

	ss.TraversAll(arr2del,fdo)

	for _,sb:=range arr2del.arrdel{
		ss.DelStream(sb)
	}

}

func (ss *streamstore)Run()  {
	select {
	case <-ss.tick:
		if tools.GetNowMsTime() - ssLastAccessTime > ssTimeOutTV{
			ssLastAccessTime = tools.GetNowMsTime()
			ss.doTimeOut()
		}
	}
}
