package stream

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsdht/nbserr"
	"io"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/translayer/ackmessage"
	"reflect"
)


type cacheblock struct {
	store.UdpMsg
	cnt int
	lastSendTime int64
}


type udpstream struct {
	conn netcommon.UdpConn
	apptyp uint32
	mtu int32
	timeout int32
	maxcache int32
	maxcnt int32
	curcnt int32
	resendtimetv int32
	lastAckPos uint64
	lastsendtime int64
	try2snd uint64
	udpmsgcache map[uint64]*cacheblock
	parent store.UdpMsg
	um store.UdpMsg
	ackchan chan interface{}
}


type UdpStream interface {
	ReliableSend(reader io.Reader) error
	SetAppTyp(typ uint32)
	SetMtu(mtu int32)
	SetTimeOut(tv int32)
	SetMaxCache(cache int32)
	SetResendTimeOut(tv int32)
	GetStreamId() (uint64,error)
}

var (
	udpsendstreamerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_STREAM_DEFAULT_ERR,Errmsg:"Send Stream Error"}
	udpstreamreadererr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_STREAM_READER_ERR,Errmsg:"Send Stream reader Error"}
	udpstreamconnerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_STREAM_CONN_ERR,Errmsg:"Send Stream connection Error"}
	udpstreamtimeouterr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"Send Stream Time Out Error"}
)

func NewUdpStream(conn netcommon.UdpConn,delayInit bool) UdpStream  {
	us:=&udpstream{conn:conn}
	us.mtu = 544
	us.timeout = 8000   //8 second
	us.maxcache = 16*(1<<10)
	us.resendtimetv = 1000
	us.udpmsgcache = make(map[uint64]*cacheblock)

	if !delayInit{
		us.um = store.NewUdpMsg(nil,0)
	}

	return us
}

var (
	sendnoneerr int = 0
	overmaxcnt int = 1
	readfinish int = 2
	readerr int = 3
	senderr int = 4
	sendfinish = 5
)


func sendUm(um store.UdpMsg,conn netcommon.UdpConn) error {
	
	if d2snd,err:=um.Serialize();err!=nil{
		return udpsendstreamerr
	}else{
		if err:=conn.Send(d2snd,store.UDP_STREAM); err!=nil{
			return err
		}
	}

	return nil
}

func (us *udpstream)sendBlk(reader io.Reader) int{

	for{
		if us.curcnt >= us.maxcnt{
			return overmaxcnt
		}
		buf:=make([]byte,us.mtu)
		n,err:=reader.Read(buf)
		if n>0 || (n==0 && err!=nil && err==io.EOF){
			var um store.UdpMsg
			if us.parent == nil {
				if us.um == nil {
					um = store.NewUdpMsg(buf[:n], us.apptyp)
				}else {
					um = us.um
					um.SetData(buf[:n])
					um.SetAppTyp(us.apptyp)
				}
				us.parent = um
				um.SetLastPos()
				us.ackchan=make(chan interface{},us.maxcnt+1)
				um.SetInform(&us.ackchan)
				ms:=store.GetBlockStoreInstance()
				ms.AddMessageWithParam(um,us.timeout,store.UDP_STREAM)
			}else{
				um = us.parent.NxtPos(buf[:n])
			}
			if err!=nil && err==io.EOF {
				um.SetLastFlag(true)
			}
			if err:=sendUm(um,us.conn);err!=nil {
				return senderr
			}
			us.curcnt ++
			us.try2snd += uint64(n)
			us.lastsendtime = tools.GetNowMsTime()

			us.udpmsgcache[um.GetPos()] = &cacheblock{um,0,tools.GetNowMsTime()}

		}
		if err!=nil{
			if err==io.EOF{
				return readfinish
			}else{
				return readerr
			}
		}
	}
}

func (us *udpstream)ReliableSend(reader io.Reader) error  {
	if us.conn == nil || !us.conn.Status(){
		return udpsendstreamerr
	}

	if us.mtu <512{
		us.mtu = 512
	}

	us.maxcnt = us.maxcache/us.mtu
	if us.maxcnt < 16 {
		us.maxcnt = 16
	}

	r:=us.sendBlk(reader)
	finishflag := false
	switch r {
	case senderr:
		return 	udpstreamconnerr
	case overmaxcnt:
		//nothing todo...
	case readfinish:
		finishflag = true
	case readerr:
		return udpstreamreadererr

	}

	ctick:=make(chan int64,8)
	ts:=tools.GetNbsTickerInstance()
	ts.RegWithTimeOut(&ctick,500)

	var ack interface{}
	for {
		ack = nil
		select {
		case ack = <-us.ackchan:
		case <-ctick:
		}
		if ack!=nil {
			switch ack.(type) {
			case int64:
				if ack == store.UDP_INFORM_TIMEOUT{
					return udpstreamtimeouterr
				}
			case ackmessage.AckMessage:

				ret:=us.doAck(ack.(ackmessage.AckMessage))
				if ret == senderr{
					return udpstreamconnerr
				}
				if ret== sendfinish{
					return nil
				}
			}
		}else{
			if r:=us.doTimeOut();r==senderr{
				return udpstreamconnerr
			}
		}
		if !finishflag{
			ret := us.sendBlk(reader)
			switch ret {
			case senderr:
				return 	udpstreamconnerr
			case overmaxcnt:
				//nothing todo...
			case readfinish:
				finishflag = true
			case readerr:
				return udpstreamreadererr
			}
		}
	}
	return nil
}

func (us *udpstream)doTimeOut()  int {

	if tools.GetNowMsTime() - us.lastsendtime < int64(us.resendtimetv){
		return sendnoneerr
	}

	listkey:=reflect.ValueOf(us.udpmsgcache).MapKeys()

	for _,key:=range listkey{
		idx:=key.Uint()
		v:=us.udpmsgcache[idx]
		curtime:=tools.GetNowMsTime()
		if curtime - v.lastSendTime > 1000 && v.cnt < 3{
			v.lastSendTime = tools.GetNowMsTime()
			v.cnt ++
			um := *v
			if err:=sendUm(um,us.conn);err!=nil {
				return senderr
			}
			us.try2snd += uint64(len(um.GetData()))
			us.lastsendtime = tools.GetNowMsTime()
		}
	}
	return sendnoneerr
}

func (us *udpstream)doAck(ack ackmessage.AckMessage) int{
	pos:=ack.GetPos()

	if us.lastAckPos < pos{
		us.lastAckPos = pos
	}
	if _,ok:=us.udpmsgcache[pos];ok{
		delete(us.udpmsgcache,pos)
		us.curcnt --
	}

	if len(us.udpmsgcache)==0{
		return sendfinish
	}

	resendpos:=ack.GetResendPos()
	for _,idx:=range resendpos{
		if v,ok:=us.udpmsgcache[idx];ok{
			um := *v
			if err:=sendUm(um,us.conn);err!=nil {
				return senderr
			}
			us.try2snd += uint64(len(um.GetData()))
		}
	}

	if len(resendpos) > 0{
		us.lastsendtime = tools.GetNowMsTime()
	}

	return sendnoneerr
}

func (us *udpstream)SetAppTyp(typ uint32)  {
	us.apptyp = typ
}

func (us *udpstream)SetMtu(mtu int32)  {
	us.mtu = mtu
}

func (us *udpstream)SetTimeOut(tv int32)  {
	us.timeout = tv
}

func (us *udpstream)SetMaxCache(cache int32)  {
	us.maxcache = cache
}

func (us *udpstream)SetResendTimeOut(tv int32)  {
	us.resendtimetv = tv
}

func (us *udpstream)GetStreamId() (uint64,error)  {
	if us.um !=nil{
		return us.um.GetSn(),nil
	}

	return 0,udpsendstreamerr
}
