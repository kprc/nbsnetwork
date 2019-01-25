package packet

import (
	"fmt"
	"github.com/kprc/nbsnetwork/pb/packet"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/common/constant"
)


type UdpPacketDataer interface {
	SetACK()
	SetDataTranser()
	SetFinishACK()
	SetTryCnt(cnt uint8)
	IncTryCnt()
	GetTryCnt() uint8
	SetPos(pos uint32)
	GetPos() uint32
	SetTotalCnt(cnt uint32)
	GetTotalCnt() uint32
	SetData(data []byte)
	GetData() []byte
	SetSerialNo(sn uint64)
	GetSerialNo() uint64
	SetRcvSN(sn uint64)
	GetRcvSN() uint64
	Serialize() ([]byte,error)
	DeSerialize(data []byte) error
	SetTyp(typ uint16)
	GetTyp() uint16
	GetLength() int32
	SetLength(len int32)
	SetTransInfo(ti []byte)
	GetTransInfo() []byte
	Finished()
	IsFinished() bool
	PrintAll()

}


type UDPPacketData struct {
	serialNo uint64    //for upper protocol used
	rcvSn uint64
	totalCnt uint32    //last packet will be set,other packet will be set to 0
	posNum uint32      //current packet serial number, start number is 1
	dataTyp uint16     //data type, for transfer priority,ACK or Data
	tryCnt  uint8      //try transfer times, default is set to 0
	transInfo []byte    //transfer layer infomation,msgid stationid, head info
	finished    bool
	len   int32			//data len, for check the data len
	data  []byte
}



func NewUdpPacketData() UdpPacketDataer {
	return &UDPPacketData{}
}

func (uh *UDPPacketData)PrintAll()  {
	fmt.Println(uh.serialNo,uh.totalCnt,uh.posNum,uh.dataTyp,uh.finished,uh.len)
}


func (uh *UDPPacketData)Serialize() ([]byte,error)  {
	p := packet.UDPPacketData{}
	p.SerialNo = uh.serialNo
	p.RcvSn = uh.rcvSn
	p.TotalCnt = uh.totalCnt
	p.PosNum = uh.posNum
	p.DataType = uint32(uh.dataTyp)
	p.TryCnt = uint32(uh.tryCnt)
	p.Data = uh.data
	p.Len = uh.len
	p.TransInfo = uh.transInfo
	p.Finished = uh.finished


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
		uh.len = p.Len
		uh.transInfo = p.TransInfo
		uh.finished = p.Finished
		uh.rcvSn = p.RcvSn
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

func (uh *UDPPacketData)IncTryCnt()  {
	uh.tryCnt ++
}

func (uh *UDPPacketData)GetTryCnt()  uint8 {
	return uh.tryCnt
}

func (uh *UDPPacketData)SetTryCnt(cnt uint8)  {
	uh.tryCnt = cnt
}


func (uh *UDPPacketData)SetTyp(typ uint16)  {
	uh.dataTyp = typ
}

func (uh *UDPPacketData)GetTyp() uint16  {
	return uh.dataTyp
}


func (uh *UDPPacketData)SetACK()  {
	uh.SetTyp(constant.ACK)
}

func (uh *UDPPacketData)SetDataTranser()  {
	uh.SetTyp(constant.DATA_TRANSER)
}

func (uh *UDPPacketData)SetFinishACK()  {
	uh.SetTyp(constant.FINISH_ACK)
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

func (uh *UDPPacketData)GetLength() int32  {
	return  uh.len
}

func (uh *UDPPacketData)SetLength(len int32){
	uh.len = len
}

func (uh *UDPPacketData)SetSerialNo(sn uint64)  {
	uh.serialNo = sn
}

func (uh *UDPPacketData)SetRcvSN(sn uint64) {
	uh.rcvSn = sn
}

func (uh *UDPPacketData)GetRcvSN() uint64  {
	return uh.rcvSn
}

func (uh *UDPPacketData)SetTransInfo(ti []byte)  {
	uh.transInfo = ti
}

func (uh *UDPPacketData)GetTransInfo() []byte {
	return uh.transInfo
}

func (uh *UDPPacketData)Finished()  {
	uh.finished = true
}

func (uh *UDPPacketData)IsFinished() bool  {
	return uh.finished
}


