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
		data :=store.GetBlkAndRefresh(blk)

		um:=data.(store.UdpMsg)

		um.Inform(store.UDP_INFORM_ACK)

		return blk,nil
	}

	ms:=store.GetBlockStoreInstance()
	if v,err:=ms.FindBlockDo(ack,nil,fdo); err!=nil{
		return err
	}else {
		if v!=nil {
			ms.DelBlock(v)
		}
	}

	return nil

}
