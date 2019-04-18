package message

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/ackmessage"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/applayer"
	"fmt"
)

var (
	dataerr=nbserr.NbsErr{Errmsg:"No Data"}
)

func Recv(rblk netcommon.RcvBlock)error  {
	data:= rblk.GetConnPacket().GetData()

	um:=store.NewUdpMsg(nil,0)

	if err:=um.DeSerialize(data);err!=nil {
		return err
	}
	cb:=applayer.NewCtrlBlk(rblk,um)

	apptyp:=um.GetAppTyp()

	abs:=applayer.GetAppBlockStore()
	if _,err:=abs.Do(apptyp,cb,nil);err!=nil{
		fmt.Println(um.GetAppTyp(),um.GetSn(),um.GetRealPos(),string(um.GetData()))
	}

	//send ack
	ack:=ackmessage.GetAckMessage(um.GetSn(),um.GetPos())

	if d2snd,err := ack.Serialize();err==nil{
		rblk.GetUdpConn().Send(d2snd,store.UDP_ACK)
	}

	return nil
}
