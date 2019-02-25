package netcommon

import (
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/pb/packet"
)

var (
	CONN_PACKET_TYP_ACK uint32 = 1
	CONN_PACKET_TYP_DATA uint32 = 2
)

type connpacket struct {
	typ uint32
	data []byte
}

type ConnPacket interface {
	SetTyp(typ uint32)
	GetTyp() uint32
	SetData(data []byte)
	GetData() []byte
	Serialize() ([]byte,error)
	UnSerilize(data []byte) error
}

func NewConnPacket() ConnPacket {
	return &connpacket{}
}



func (cp *connpacket)SetTyp(typ uint32)  {
	cp.typ = typ
}

func (cp *connpacket)GetTyp() uint32  {
	return cp.typ
}

func (cp *connpacket)SetData(data []byte)  {
	cp.data = data
}

func (cp *connpacket)GetData() []byte  {
	return cp.data
}

func (cp *connpacket)Serialize() ([]byte,error)  {
	p:=&packet.UdpConnMsg{}

	p.Typ = cp.typ
	p.Data = cp.data

	return  proto.Marshal(p)
}

func (cp *connpacket)UnSerilize(data []byte) error {
	p := &packet.UdpConnMsg{}

	if err:=proto.Unmarshal(data,p); err!=nil{
		return err
	}

	cp.typ = p.Typ
	cp.data = p.Data

	return nil

}
