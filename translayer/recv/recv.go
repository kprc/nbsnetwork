package recv

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/translayer/ackmessage"
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/translayer/stream"
)

func ReceiveFromUdpConn() error  {
	cs:=netcommon.GetConnStoreInstance()

	cs.RSync()
	defer cs.RDone()

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
			stream.Recv(rcvblk)
		case store.UDP_RCV_QUIT:
			return nil
		}
	}

}

func ReceiveFromUdpStop()  {
	cs:=netcommon.GetConnStoreInstance()
	cp:=netcommon.NewConnPacket()
	cp.SetMsgTyp(store.UDP_RCV_QUIT)
	rcvblk:=netcommon.NewRcvBlock(cp,nil)

	cs.Push(rcvblk)

	cs.RWait()

}