package recv

import "github.com/kprc/nbsnetwork/netcommon"



func ReceiveFromUdpConn() error  {
	cs:=netcommon.GetConnStoreInstance()

	for{
		rcvblk:=cs.Read()
		translayerdata := rcvblk.GetConnPacket()
		uid:=translayerdata.GetUid()
		data:=translayerdata.GetData()
		msgtyp := translayerdata.GetMsgTyp()

		//send to receive module


	}
}
