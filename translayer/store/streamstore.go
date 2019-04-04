package store

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"sync"
	"hash/fnv"
	"github.com/kprc/nbsnetwork/tools"
)

type udpstreamkey struct {
	uid string
	sn uint64
}

type UdpStreamKey interface{
	SetUid(uid string)
	GetUid() string
	SetSn(sn uint64)
	GetSn() uint64
	Equals(key UdpStreamKey) bool
}

func NewUdpStreamKey() UdpStreamKey  {
	return &udpstreamkey{}
}

func NewUdpStreamKeyWithParam(uid string, sn uint64) UdpStreamKey  {
	return &udpstreamkey{uid:uid,sn:sn}
}

func (usk *udpstreamkey)SetUid(uid string)  {
	usk.uid = uid
}

func (usk *udpstreamkey)GetUid() string  {
	return usk.uid
}

func (usk *udpstreamkey)SetSn(sn uint64)  {
	usk.sn = sn
}

func (usk *udpstreamkey)GetSn() uint64  {
	return usk.sn
}

func (usk *udpstreamkey)Equals(key UdpStreamKey) bool  {
	if usk.sn == key.GetSn() && usk.uid == key.GetUid() {
		return true
	}

	return false
}

type streamblk struct {
	key UdpStreamKey
	blk interface{}
	timeoutInterval int32
	lastAccessTime int64
}


type StreamKeyInter interface {
	GetKey() UdpStreamKey
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

var fsshash = func(v interface{}) uint {
	sk:=v.(StreamKeyInter)
	s := fnv.New64()
	s.Write([]byte(sk.GetKey().GetUid()))
	h:=s.Sum64()

	h = sk.GetKey().GetSn() << 8 | h

	h = h & 0x7F

	return uint(h)

}

var fssequals = func(v1 interface{}, v2 interface{}) int {
	sk1:=v1.(StreamKeyInter)
	sk2:=v2.(StreamKeyInter)

	if sk1.GetKey().GetUid() == sk2.GetKey().GetUid() {
		if sk1.GetKey().GetSn() == sk2.GetKey().GetSn(){
			return 0
		}
	}

	return 1
}

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

	hl:=hashlist.NewHashList(0x80,fsshash,fssequals)
	ss:=&streamstore{hl,make(chan int64,64)}

	t:=tools.GetNbsTickerInstance()
	t.Reg(&ss.tick)

	ssLastAccessTime = tools.GetNowMsTime()

	return ss
}

func (ss *streamstore)addStream(v interface{}, timeinterval int32)  {
	sb:=&streamblk{blk:v}
	sb.key = v.(StreamKeyInter).GetKey()
	sb.timeoutInterval = timeinterval
	sb.lastAccessTime = tools.GetNowMsTime()

	ss.Add(sb)
}

func (ss *streamstore)AddStream(v interface{})  {
	ss.addStream(v,5000)
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














