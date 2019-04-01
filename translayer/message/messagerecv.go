package message

import (
	"github.com/kprc/nbsnetwork/pb/udpmessage"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/gogo/protobuf/proto"
	"fmt"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/ackmessage"
	"github.com/kprc/nbsnetwork/translayer/store"
)

var (
	dataerr=nbserr.NbsErr{Errmsg:"No Data"}
)

func Recv(rblk netcommon.RcvBlock)error  {
	data:= rblk.GetConnPacket().GetData()
	um:=&udpmessage.Udpmsg{}

	if err:=proto.Unmarshal(data,um);err!=nil {
		return err
	}

	//send ack
	ack:=ackmessage.GetAckMessage(um.GetSn(),um.GetPos())

	if d2snd,err := ack.Serialize();err!=nil{
		rblk.GetUdpConn().Send(d2snd,store.UDP_ACK)
	}

	fmt.Println(um.Sn,um.Pos)

	return nil
}
