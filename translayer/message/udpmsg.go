package message

import (
	"github.com/gogo/protobuf/proto"
	pbmsg "github.com/kprc/nbsnetwork/pb/message"
)

type udpmsg struct {
	sn uint64
	msgtyp int32
	data []byte
}

type UdpMsg interface {
	SetSn(sn uint64)
	GetSn() uint64
	SetMsgTyp(msgtyp int32)
	GetMsgTyp() int32
	SetData(data []byte)
	GetData() []byte
	Serialize() ([]byte,error)
	DeSerialize(data []byte) error
}

func (um *udpmsg)SetSn(sn uint64)  {
	um.sn = sn
}
func (um *udpmsg)GetSn() uint64  {
	return um.sn
}

func (um *udpmsg)SetMsgTyp(msgtyp int32){
	um.msgtyp = msgtyp
}
func (um *udpmsg)GetMsgTyp() int32{
	return um.msgtyp
}
func (um *udpmsg)SetData(data []byte){
	um.data = data
}
func (um *udpmsg)GetData() []byte{
	return um.data
}
func (um *udpmsg)Serialize() ([]byte,error){
	pbum := &pbmsg.Udpmsg{}

	pbum.Data= um.data
	pbum.Msgtyp = um.msgtyp
	pbum.Sn = um.sn
	d,err:=proto.Marshal(pbum);
	if err!=nil{
		return nil,err
	}

	return d,nil

}

func (um *udpmsg)DeSerialize(data []byte) error{
	pbum:=&pbmsg.Udpmsg{}

	if err:=proto.UnmarshalMerge(data,pbum); err!=nil{
		return err
	}

	return nil

}
