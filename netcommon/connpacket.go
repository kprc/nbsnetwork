package netcommon

import (
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/pb/packet"
)

var (
	CONN_PACKET_TYP_KA uint32 = 1
	CONN_PACKET_TYP_DATA uint32 = 2
)

type connpacket struct {
	typ uint32    //ka packet or data packet
	uid []byte
	data []byte
}

type ConnPacket interface {
	SetUid(uid []byte)
	GetUid() []byte
	SetTyp(typ uint32)
	GetTyp() uint32
	SetMsgTyp(typ uint32)
	GetMsgTyp() uint32
	SetData(data []byte)
	GetData() []byte
	Serialize() ([]byte,error)
	DeSerialize(data []byte) error
}

func NewConnPacket() ConnPacket {
	return &connpacket{}
}

func (cp *connpacket)SetUid(uid []byte)  {
	cp.uid = uid
}

func (cp *connpacket)GetUid() []byte  {
	return cp.uid
}

func (cp *connpacket)SetTyp(typ uint32)  {
	var typ1,typ2 uint32
	typ1 = cp.typ & 0x00FFFFFF
	typ2 = ((typ << 24) & 0xFF000000) | typ1
	cp.typ = typ2
}

func (cp *connpacket)GetTyp() uint32  {
	var typ1 uint32 = cp.typ

	typ1 = (typ1 >> 24) & 0x000000FF

	return typ1
}

func (cp *connpacket)SetMsgTyp(typ uint32)  {
	var typ1,typ2 uint32
	typ1 = cp.typ & 0xFF000000
	typ2 = typ & 0x00FFFFFF

	cp.typ = typ1 | typ2
}

func (cp *connpacket)GetMsgTyp() uint32  {
	return cp.typ & 0x00FFFFFF
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
	p.Uid = cp.uid

	return  proto.Marshal(p)
}

func (cp *connpacket)DeSerialize(data []byte) error {
	p := &packet.UdpConnMsg{}

	if err:=proto.Unmarshal(data,p); err!=nil{
		return err
	}

	cp.typ = p.Typ
	cp.data = p.Data
	cp.uid = p.Uid

	return nil

}
