package message

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/constant"
)

type reliablemsg struct {
	conn netcommon.UdpConn
	apptyp uint32
	timeout int32
	um store.UdpMsg
}


type ReliableMsg interface {
	ReliableSend(data []byte) (err error)
	SetAppTyp(typ uint32)
	SetTimeOut(timeout int32)
	GetSn() uint64
}

var(
	udpsendtimeouterr=nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"Udp Send Timeout"}
	udpsenddefaulterr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_DEFAULT_ERR,Errmsg:"Send Message Error"}
	udpsendoutimeserr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_OUT_TIMES,Errmsg:"Udp Send too much times"}
)

func NewReliableMsg(conn netcommon.UdpConn) ReliableMsg {
	return &reliablemsg{conn:conn,timeout:int32(constant.UDP_MESSAGE_STORE_TIMEOUT)}
}

func NewReliableMsgWithDelayInit(conn netcommon.UdpConn, delayInit bool) ReliableMsg {
	rm := &reliablemsg{conn:conn,timeout:int32(constant.UDP_MESSAGE_STORE_TIMEOUT)}
	if !delayInit {
		um:=store.NewUdpMsg(nil,0)
		rm.um = um
	}

	return rm
}


func sendUm(um store.UdpMsg,conn netcommon.UdpConn) error {

	if d2snd,err:=um.Serialize();err!=nil{
		return udpsenddefaulterr
	}else{
		if err:=conn.Send(d2snd,store.UDP_MESSAGE); err!=nil{
			return err
		}
	}

	return nil
}


func (rm *reliablemsg) ReliableSend(data []byte) (err error) {

	if rm.conn == nil || !rm.conn.Status() || len(data) == 0 || len(data) > 1024{
		return udpsenddefaulterr
	}
	var um store.UdpMsg
	if rm.um !=nil{
		um = rm.um
		um.SetData(data)
		um.SetAppTyp(rm.apptyp)
	}else {
		um = store.NewUdpMsg(data, rm.apptyp)
	}

	c:=make(chan interface{},1)
	um.SetInform(&c)


	ms:=store.GetBlockStoreInstance()

	ms.AddMessageWithParam(um,rm.timeout,store.UDP_MESSAGE)

	if err:=sendUm(um,rm.conn);err!=nil {
		return err
	}

	rc:=<-c
	r:=rc.(int64)

	if r == store.UDP_INFORM_ACK{
		ms.DelMessage(um)
		return nil
	}else if r == store.UDP_INFORM_TIMEOUT{
		return udpsendtimeouterr
	}

	return udpsenddefaulterr
}

func (rm *reliablemsg)SetAppTyp(typ uint32)  {
	rm.apptyp = typ
}

func (rm *reliablemsg)SetTimeOut(timeout int32)  {
	rm.timeout = timeout
}

func (rm *reliablemsg)GetSn() uint64  {
	if rm.um != nil{
		return rm.um.GetSn()
	}

	return 0
}
