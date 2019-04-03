package message

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsdht/nbserr"
)

type reliablemsg struct {
	conn netcommon.UdpConn
	timeout int32
}


type ReliableMsg interface {
	ReliableSend(data []byte) (err error)
	SetTimeOut(timeout int32)

}

var(
	udpsendtimeouterr=nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"Udp Send Timeout"}
	udpsenddefaulterr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_DEFAULT_ERR,Errmsg:"Send Message Error"}
	udpsendoutimeserr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_OUT_TIMES,Errmsg:"Udp Send too much times"}
)

func NewReliableMsg(conn netcommon.UdpConn) ReliableMsg {
	return &reliablemsg{conn:conn}
}

func SendUm(um store.UdpMsg,conn netcommon.UdpConn) error {
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

	if rm.conn == nil || !rm.conn.Status(){
		return udpsenddefaulterr
	}

	um:=store.NewUdpMsg(data)

	c:=make(chan interface{},1)
	um.SetInform(&c)

	ms:=store.GetBlockStoreInstance()

	ms.AddMessageWithParam(um,rm.timeout)

	if err:=SendUm(um,rm.conn);err!=nil {
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


func (rm *reliablemsg)SetTimeOut(timeout int32)  {
	rm.timeout = timeout
}

