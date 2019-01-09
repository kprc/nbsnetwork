package packet

import (
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/pb"
)

type udpack struct {
	resend []uint32
	rcved uint32

}

type udpresult struct {
	serialNo uint64
	finished bool
	udpack
}

type UdpAcker interface {
	AppendResend(ids...uint32)
	SetRcved(id uint32)
	GetReSend() []uint32
	GetRcved() uint32
}

type UdpResulter interface {
	SetSerialNo(sn uint64)
	GetSerialNo() uint64
	Finished()
	IsFinished() bool
	Serialize() ([]byte,error)
	DeSerialize(buf []byte) error
	UdpAcker
}

func NewUdpResult(sn uint64) UdpResulter {
	ur := &udpresult{serialNo:sn}
	ur.resend = make([]uint32,0)

	return ur
}

func (ur *udpresult)Finished() {
	ur.finished = true
}

func (ur *udpresult)IsFinished() bool  {
	return ur.finished
}


func (ur *udpresult) SetSerialNo(sn uint64) {
	ur.serialNo = sn
}

func (ur *udpresult) GetSerialNo() uint64 {
	return ur.serialNo
}

func (ur *udpresult)Serialize() ([]byte,error)  {
	ua:=&packet.UdpAck{}
	ua.SerialNo = ur.serialNo
	ua.Ack = ur.rcved
	ua.Resend = ur.resend
	ua.Finished = ur.finished

	return proto.Marshal(ua)
}

func (ur *udpresult)DeSerialize(buf []byte) error{
	ua:=&packet.UdpAck{}

	err:=proto.Unmarshal(buf,ua)

	if err==nil {
		ur.resend = ua.Resend
		ur.rcved = ua.Ack
		ur.serialNo = ua.SerialNo
		ur.finished = ua.Finished
	}

	return err
}

func (ua *udpack) AppendResend(ids...uint32) {
	for _,id:= range ids{
		ua.resend = append(ua.resend,id)
	}

}


func (ua *udpack) SetRcved(id uint32) {
	ua.rcved = id
}


func (ua *udpack) GetReSend() []uint32 {
	return  ua.resend
}


func (ua *udpack) GetRcved() uint32 {
	return ua.rcved
}




