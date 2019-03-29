package ackmessage


type Ackid struct {
	sn uint64
	pos uint64
}

type ackmessage struct {
	Ackid
	resendids []*Ackid
}


type AckMessage interface {
	GetSN() uint64
	GetPos() uint64
	Append(sn,pos uint64)
	GetResendID() []*Ackid
}

func (aid *Ackid)GetSN() uint64 {
	return aid.sn
}

func (aid *Ackid)GetPos() uint64  {
	return aid.pos
}

func NewAckMessage(sn,pos uint64) AckMessage {
	return &ackmessage{Ackid{sn,pos},make([]*Ackid,0)}
}

func (am *ackmessage)Append(sn,pos uint64)  {
	aid:=&Ackid{sn,pos}

	am.resendids = append(am.resendids,aid)
}

func (am *ackmessage)GetResendID() []*Ackid  {
	return am.resendids
}





