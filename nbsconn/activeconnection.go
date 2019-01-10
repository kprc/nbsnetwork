package nbsconn

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/netcommon"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type udpRW struct {
	urw netcommon.UdpReaderWriterer
	lastAccessTime int64
	refcnt int32
}

type mapRW struct {
	lock sync.RWMutex
	rw map[string]*udpRW
}

func (mr *mapRW)Count() int  {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	return len(mr.rw)
}

func (mr *mapRW)TimeOut()  {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	now:=time.Now().Unix()

	keys:=reflect.ValueOf(mr.rw).MapKeys()
	
	for _,key:=range keys{
		urw,_ := mr.rw[key.String()]
		if (now - urw.lastAccessTime) > constant.UDP_CONN_TIMEOUT {
			mr.delRW(urw.urw)
		}
	}
}


func (mr *mapRW)AddMRW(urw netcommon.UdpReaderWriterer) error {
	mr.lock.Lock()
	defer mr.lock.Unlock()
	if _,ok:=mr.rw[urw.AddrString()];ok {
		return nbserr.NbsErr{Errmsg:"connection is exists",ErrId:nbserr.UDP_CONN_EXISTS}
	}

	u := &udpRW{urw:urw,lastAccessTime:time.Now().Unix()}

	atomic.AddInt32(&u.refcnt,1)

	mr.rw[urw.AddrString()] = u

	return nil

}

func (mr *mapRW)delRW(urw netcommon.UdpReaderWriterer)  {
	if u,ok:=mr.rw[urw.AddrString()];ok{
		atomic.AddInt32(&u.refcnt,-1)
		if atomic.LoadInt32(&u.refcnt) <= 0 {
			delete(mr.rw,urw.AddrString())
		}
	}
}

func (mr *mapRW)DelRW(urw netcommon.UdpReaderWriterer)  {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	mr.delRW(urw)

}

func (mr *mapRW)GetRw() netcommon.UdpReaderWriterer  {

	mr.lock.RLock()
	defer mr.lock.RUnlock()
	keys := reflect.ValueOf(mr.rw).MapKeys()

	if len(keys) >0 {
		u:=mr.rw[keys[0].String()]
		u.lastAccessTime = time.Now().Unix()
		atomic.AddInt32(&u.refcnt,1)

		return u.urw
	}

	return nil

}

func (mr *mapRW)PutRW(urw netcommon.UdpReaderWriterer)  {
	mr.lock.RLock()


	if u,ok:=mr.rw[urw.AddrString()];ok{
		u.lastAccessTime = time.Now().Unix()
		atomic.AddInt32(&u.refcnt,-1)

		if atomic.LoadInt32(&u.refcnt) <= 0 {
			mr.lock.RUnlock()
			mr.lock.Lock()
			delete(mr.rw,urw.AddrString())
			mr.lock.Unlock()

		}else {
			mr.lock.RUnlock()
		}

	}else {
		mr.lock.RUnlock()
	}

}


type mapConn struct {
	lock sync.Mutex
	collectConn map[string]*mapRW
}



var (
	mconnInstance MapConn
	glock sync.Mutex
)

type MapConn interface {
	AddConn(remoteStationId string, urw netcommon.UdpReaderWriterer) error
	DelConn(remoteStationId string, urw netcommon.UdpReaderWriterer)
	GetActiveConn(remoteStationId string) netcommon.UdpReaderWriterer
	PutActiveConn(remoteStationId string, urw netcommon.UdpReaderWriterer)
	TimeOut()
}



func GetMapConnInstance() MapConn  {
	if mconnInstance == nil{
		glock.Lock()
		if mconnInstance == nil{
			mconnInstance = &mapConn{collectConn:make(map[string]*mapRW)}
		}
		glock.Unlock()
	}

	return mconnInstance
}

func (mc *mapConn) AddConn(remoteStationId string, urw netcommon.UdpReaderWriterer) error {

	mc.lock.Lock()
	defer mc.lock.Unlock()

	 mrw,ok :=mc.collectConn[remoteStationId]
	 if !ok{
	 	mrw = &mapRW{rw:make(map[string]*udpRW)}
	 }

	err:=mrw.AddMRW(urw)
	if err!=nil {
		return err
	}

	mc.collectConn[remoteStationId] = mrw

	return nil
}

func (mc *mapConn)DelConn(remoteStationId string, urw netcommon.UdpReaderWriterer)  {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mrw,ok :=mc.collectConn[remoteStationId]
	if !ok{
		return
	}

	mrw.DelRW(urw)

	return

}

func (mc *mapConn)GetActiveConn(remoteStationId string) netcommon.UdpReaderWriterer  {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mrw,ok :=mc.collectConn[remoteStationId]
	if !ok{
		return mrw.GetRw()
	}

	return nil
}

func (mc *mapConn)PutActiveConn(remoteStationId string, urw netcommon.UdpReaderWriterer)  {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mrw,ok:=mc.collectConn[remoteStationId]
	if !ok {
		return
	}

	mrw.PutRW(urw)

}

func (mc *mapConn)timeOut()  {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	keys:=reflect.ValueOf(mc.collectConn).MapKeys()

	for _,key :=range keys{

		mr:=mc.collectConn[key.String()]

		mr.TimeOut()

		if mr.Count() ==0 {
			delete(mc.collectConn,key.String())
		}
	}
}

func (mc *mapConn)TimeOut()  {
	for{
		time.Sleep(time.Second*60)
		mc.timeOut()
	}

}


