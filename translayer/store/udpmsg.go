package store

import (
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/pb/udpmessage"
	"sync/atomic"
)

var gSerialNumber uint64

const(
	UDP_MESSAGE uint32=0
	UDP_STREAM uint32=1
	UDP_ACK uint32 = 2

	UDP_INFORM_ACK int64 = 0
	UDP_INFORM_TIMEOUT int64 = 1
)

type udpmsg struct {
	sn uint64
	pos uint64
	data []byte
	inform *chan interface{}
}

type UdpMsg interface {
	SetSn(sn uint64)
	GetSn() uint64
	SetPos(pos uint64)
	GetPos() uint64
	SetData(data []byte)
	GetData() []byte
	Serialize() ([]byte,error)
	DeSerialize(data []byte) error
	NxtPos(data []byte) UdpMsg
	SetInform(c *chan interface{})
	Inform(typ int64)
}

func getNextSerialNum() uint64 {
	return atomic.AddUint64(&gSerialNumber,1)
}

func NewUdpMsg(data []byte) UdpMsg  {
	um:=&udpmsg{sn:getNextSerialNum(),data:data}

	return um
}



func (um *udpmsg)SetInform(c *chan interface{}) {
	um.inform = c
}

func (um *udpmsg)Inform(typ int64)  {
	if um.inform != nil {
		select {
		case *um.inform <- typ:
		default:
			//nothing to do
		}
	}
}



func (um *udpmsg)NxtPos(data []byte) UdpMsg  {
	um1 := &udpmsg{}
	um1.sn = um.sn
	um1.data = data
	um1.pos ++

	return um1
}





func (um *udpmsg)SetSn(sn uint64)  {
	um.sn = sn
}
func (um *udpmsg)GetSn() uint64  {
	return um.sn
}

func (um *udpmsg)SetPos(pos uint64){
	um.pos =pos
}
func (um *udpmsg)GetPos() uint64{
	return um.pos
}

func (um *udpmsg)SetData(data []byte){
	um.data = data
}
func (um *udpmsg)GetData() []byte{
	return um.data
}
func (um *udpmsg)Serialize() ([]byte,error){
	pbum := &udpmessage.Udpmsg{}

	pbum.Data= um.data
	pbum.Sn = um.sn
	pbum.Pos = um.pos
	d,err:=proto.Marshal(pbum);
	if err!=nil{
		return nil,err
	}

	return d,nil

}

func (um *udpmsg)DeSerialize(data []byte) error{
	pbum:=&udpmessage.Udpmsg{}

	if err:=proto.UnmarshalMerge(data,pbum); err!=nil{
		return err
	}

	um.sn = pbum.Sn
	um.pos = pbum.Pos
	um.data = pbum.Data

	return nil
}

