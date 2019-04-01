package recv

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/translayer/ackmessage"
	"github.com/kprc/nbsnetwork/translayer/message"
)

func ReceiveFromUdpConn() error  {
	cs:=netcommon.GetConnStoreInstance()

	for{
		rcvblk:=cs.Read()
		translayerdata := rcvblk.GetConnPacket()

		msgtyp := translayerdata.GetMsgTyp()
		switch msgtyp {
		case store.UDP_ACK:
			ackmessage.AckRecv(rcvblk)
		case store.UDP_MESSAGE:
			message.Recv(rcvblk)
		case store.UDP_STREAM:
			//TO DO...
		}
	}
}
