package ackmessage

import (
	"github.com/kprc/nbsnetwork/pb/udpmessage"
	"github.com/gogo/protobuf/proto"

)

type Ackid struct {
	sn uint64
	pos uint64
}

type ackmessage struct {
	Ackid
	resendpos []uint64
}


type AckMessage interface {
	GetSn() uint64
	GetPos() uint64
	SetSn(sn uint64)
	SetPos(pos uint64)
	Append(pos uint64)
	GetResendPos() []uint64
	SetResendPos(arr []uint64)
	Serialize() ([]byte,error)
	Deserialize(data []byte) error
}

func (aid *Ackid)GetSn() uint64 {
	return aid.sn
}

func (aid *Ackid)GetPos() uint64  {
	return aid.pos
}

func (aid *Ackid)SetSn(sn uint64)  {
	aid.sn = sn
}

func (aid *Ackid)SetPos(pos uint64)  {
	aid.pos = pos
}


func NewAckMessage(sn,pos uint64) AckMessage {
	return &ackmessage{Ackid{sn,pos},make([]uint64,0)}
}

func (am *ackmessage)Append(pos uint64)  {
	//aid:=&Ackid{sn,pos}

	am.resendpos = append(am.resendpos,pos)
}

func (am *ackmessage)GetResendPos() []uint64  {
	return am.resendpos
}

func (am *ackmessage)SetResendPos(arr []uint64)  {
	am.resendpos = arr
}

func (am *ackmessage)Serialize() ([]byte,error)  {
	uma:=&udpmessage.Udpmsgack{}
	umid:=&udpmessage.Udpmsgid{Sn:am.GetSn(),Pos:am.GetPos()}

	uma.Uid = umid

	arrpos := make([]uint64,0)

	for _,aid:=range am.resendpos{
		arrpos=append(arrpos,aid)
	}

	uma.Arrpos = arrpos

	return proto.Marshal(uma)
}

func (am *ackmessage)Deserialize(data []byte) error {
	uma:=&udpmessage.Udpmsgack{}

	if err:= proto.Unmarshal(data,uma);err!=nil{
		return err
	}

	am.SetSn(uma.Uid.GetSn())
	am.SetPos(uma.Uid.GetPos())
	am.resendpos = uma.Arrpos

	return nil

}


