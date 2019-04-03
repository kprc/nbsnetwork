package ackmessage

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/store"
)

func AckRecv(rblk netcommon.RcvBlock)  error{
	data:=rblk.GetConnPacket().GetData()
	ack:=&ackmessage{}

	if err:=ack.Deserialize(data);err!=nil{
		return err
	}

	fdo:= func(arg interface{},blk interface{}) (v interface{},err error) {

		typ := store.GetMsgTyp(blk)

		if typ == store.UDP_MESSAGE {
			data :=store.GetBlk(blk,true)

			um:=data.(store.UdpMsg)

			um.Inform(store.UDP_INFORM_ACK)

		}else if typ == store.UDP_STREAM{
			data :=store.GetBlk(blk,false)

			um:=data.(store.UdpMsg)

			um.Inform(arg)
		}
		return blk,nil
	}

	ms:=store.GetBlockStoreInstance()
	if _,err:=ms.FindMessageDo(ack,nil,fdo); err!=nil{
		return err
	}else {
		//if v!=nil {
		//	ms.DelMessage(v)
		//}
		//nothing todo...
	}

	return nil

}
