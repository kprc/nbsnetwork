package nbsnetwork

import (
	"sync/atomic"
)


type UdpPacketDataer interface {
	SetPing()
	SetACK()
	SetDataTranser()
	SetTryCnt(cnt uint8)
	GetTryCnt() uint8
	SetPos(pos uint32)
	GetPos() uint32
	SetTotalCnt(cnt uint32)
	GetTotalCnt() uint32
	SetData(data []byte)
	GetData() []byte
	GetSerialNo() uint64
}


type UDPPacketData struct {
	serialNo uint64    //for upper protocol used
	totalCnt uint32    //last packet will be set,other packet will be set to 0
	posNum uint32      //current packet serial number
	dataTyp uint16     //data type, for transfer priority
	tryCnt  uint8      //retry transfer times
	pad8    uint8
	pad32   uint32
	data  []byte
}

func NewUdpPacketData(sn uint64, dataType uint16) UdpPacketDataer {
	upd := &UDPPacketData{serialNo:sn,dataTyp:dataType}
	return upd
}



func (uh *UDPPacketData)SetTotalCnt(cnt uint32)  {
	uh.totalCnt = cnt
}

func (uh *UDPPacketData)GetTotalCnt() uint32  {
	return uh.totalCnt
}

func (uh *UDPPacketData)SetPos(pos uint32)  {
	atomic.StoreUint32(&uh.posNum,pos)
}

func (uh *UDPPacketData)GetPos() uint32  {
	return atomic.LoadUint32(&uh.posNum)
}

func (uh *UDPPacketData)GetTryCnt()  uint8 {
	return uh.tryCnt
}

func (uh *UDPPacketData)SetTryCnt(cnt uint8)  {
	uh.tryCnt = cnt
}

func (uh *UDPPacketData)setTyp(typ uint16)  {
	uh.dataTyp = typ
}

func (uh *UDPPacketData)GetTyp() uint16  {
	return uh.dataTyp
}

func (uh *UDPPacketData)SetPing()  {
	uh.setTyp(PING)
}

func (uh *UDPPacketData)SetACK()  {
	uh.setTyp(ACK)
}

func (uh *UDPPacketData)SetDataTranser()  {
	uh.setTyp(DATA_TRANSER)
}

func (uh *UDPPacketData)SetData(data []byte)  {
	uh.data = data
}

func (uh *UDPPacketData)GetData() []byte  {
	return uh.data
}

func (uh *UDPPacketData)GetSerialNo() uint64  {
	return uh.serialNo
}