package nbsnetwork

import (
	"sync/atomic"
	"github.com/kprc/nbsnetwork/pb"
	"github.com/gogo/protobuf/proto"
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
	Serialize() ([]byte,error)
	DeSerialize(data []byte) error
}


type UDPPacketData struct {
	serialNo uint64    //for upper protocol used
	totalCnt uint32    //last packet will be set,other packet will be set to 0
	posNum uint32      //current packet serial number
	dataTyp uint16     //data type, for transfer priority
	tryCnt  uint8      //try transfer times
	pad8    uint8
	pad32   uint32
	data  []byte
}

func NewUdpPacketData(sn uint64, dataType uint16) UdpPacketDataer {
	upd := &UDPPacketData{serialNo:sn,dataTyp:dataType}
	return upd
}


//in the future i will Serialize and DeSerialize by myself
//func (uh *UDPPacketData)Serialize() []byte  {
//	buf := new(bytes.Buffer)
//	binary.Write(buf,binary.BigEndian,uh.serialNo)
//	binary.Write(buf,binary.BigEndian,uh.totalCnt)
//	binary.Write(buf,binary.BigEndian,uh.posNum)
//	binary.Write(buf,binary.BigEndian,uh.dataTyp)
//	binary.Write(buf,binary.BigEndian,uh.tryCnt)
//	buf.Write(uh.data)
//
//	return buf.Bytes()
//}
//
//func (uh *UDPPacketData)DeSerialize(data []byte) {
//	buf := new(bytes.Buffer)
//	buf.Write(data)
//	binary.Read(buf,binary.BigEndian,&uh.serialNo)
//	binary.Read(buf,binary.BigEndian,&uh.totalCnt)
//	binary.Read(buf,binary.BigEndian,&uh.posNum)
//	binary.Read(buf,binary.BigEndian,&uh.dataTyp)
//	binary.Read(buf,binary.BigEndian,&uh.tryCnt)
//
//	uh.data = buf.Bytes()
//}


func (uh *UDPPacketData)Serialize() ([]byte,error)  {
	p := packet.UDPPacketData{}
	p.SerialNo = uh.serialNo
	p.TotalCnt = uh.totalCnt
	p.PosNum = uh.posNum
	p.DataType = uint32(uh.dataTyp)
	p.TryCnt = uint32(uh.tryCnt)
	p.Data = uh.data

	return proto.Marshal(&p)
}

func (uh *UDPPacketData)DeSerialize(data []byte) error {
	p := packet.UDPPacketData{}
	err := proto.Unmarshal(data,&p)
	if  err == nil{
		uh.serialNo = p.SerialNo
		uh.totalCnt = p.TotalCnt
		uh.posNum = p.PosNum
		uh.dataTyp = uint16(p.DataType)
		uh.tryCnt = uint8(p.TryCnt)
		uh.data = p.Data
	}

	return err
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