package message

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsdht/nbserr"
)

type reliablemsg struct {
	conn netcommon.UdpConn
	uid string
	timeout int32
	resendtimes int32
	step int32
}


type ReliableMsg interface {
	ReliableSend(data []byte) (err error)
	SetTimeOut(timeout int32)
	SetReSendTimes(resendtimes int32)
	SetStep(step int32)
	GetUid() string
	//WaitAnswer()
}

var(
	udpsendtimeouterr=nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"Udp Send Timeout"}
	udpsenddefaulterr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_DEFAULT_ERR,Errmsg:"Send Error"}
)

func NewReliableMsg(conn netcommon.UdpConn,uid string) ReliableMsg {
	return &reliablemsg{conn:conn,uid:uid}
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

	c:=make(chan int64,0)
	um.SetInform(&c)

	ms:=store.GetBlockStoreInstance()

	ms.AddBlockWithParam(um,rm.uid,rm.timeout,rm.resendtimes,rm.step)

	if err:=SendUm(um,rm.conn);err!=nil {
		return err
	}

	r:=<-c

	if r == store.UDP_INFORM_ACK{
		return nil
	}else if r == store.UDP_INFORM_TIMEOUT{
		return udpsendtimeouterr
	}

	return udpsenddefaulterr

}

func (rm *reliablemsg)GetUid() string  {
	return rm.uid
}

func (rm *reliablemsg)SetTimeOut(timeout int32)  {
	rm.timeout = timeout
}

func (rm *reliablemsg)SetReSendTimes(resendtimes int32)  {
	rm.resendtimes = resendtimes
}

func (rm *reliablemsg)SetStep(step int32)  {
	rm.step = step
}