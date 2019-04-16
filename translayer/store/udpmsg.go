package store

import (
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/pb/udpmessage"
	"sync/atomic"
	"fmt"
)

var gSerialNumber uint64 = 0x20151031

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
	last bool
	data []byte
	lastpos uint64
	inform *chan interface{}
}

type UdpMsg interface {
	SetSn(sn uint64)
	GetSn() uint64
	SetLastPos()
	SetPos(pos uint64)
	GetPos() uint64
	SetRealPos(pos uint32)
	GetRealPos() uint32
	SetAppTyp(typ uint32)
	GetAppTyp() uint32
	SetLastFlag(b bool)
	GetLastFlag() bool
	SetData(data []byte)
	GetData() []byte
	Serialize() ([]byte,error)
	DeSerialize(data []byte) error
	NxtPos(data []byte) UdpMsg
	SetInform(c *chan interface{})
	Inform(v interface{})
	Print()
}

func getNextSerialNum() uint64 {
	return atomic.AddUint64(&gSerialNumber,1)
}

func NewUdpMsg(data []byte,typ uint32) UdpMsg  {
	um:=&udpmsg{sn:getNextSerialNum(),data:data}

	um.SetAppTyp(typ)

	return um
}

func (um *udpmsg)Print(){
	fmt.Println("sn:",um.sn,"pos",um.GetRealPos(),"apptype:",um.GetAppTyp(),"last:",um.last)
}

func (um *udpmsg)SetInform(c *chan interface{}) {
	um.inform = c
}

func (um *udpmsg)Inform(v interface{})  {
	if um.inform != nil {
		select {
		case *um.inform <- v:
		default:
			//nothing to do
		}
	}
}

func (um *udpmsg)NxtPos(data []byte) UdpMsg  {
	um1 := &udpmsg{}
	um1.sn = um.sn
	um1.last = um.last
	um1.data = data

	um.lastpos ++
	um1.pos = um.lastpos

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
func (um *udpmsg)SetLastPos()  {
	um.lastpos = um.pos
}

func (um *udpmsg)GetPos() uint64{
	return um.pos
}

func (um *udpmsg)SetRealPos(pos uint32){
	var l uint64

	l = uint64(pos)

	h:=um.pos >> 32

	um.pos = h | l
}


func (um *udpmsg)GetRealPos() uint32{
	return uint32(um.pos  & 0xFFFFFFFF)
}


func (um *udpmsg)SetAppTyp(typ uint32)  {
	l:=um.pos  & 0xFFFFFFFF
	var h uint64
	h = uint64(typ)
	h = h<<32

	um.pos = h | l
}

func (um *udpmsg)GetAppTyp() uint32  {
	typ := uint32(um.pos >> 32)

	return typ
}

func (um *udpmsg)SetLastFlag(b bool){
	um.last = b
}

func (um *udpmsg)GetLastFlag() bool  {
	return um.last
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
	pbum.Last = um.last
	d,err:=proto.Marshal(pbum)
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
	um.last = pbum.Last

	return nil
}

