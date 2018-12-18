package nbsnetwork


type udpack struct {
	resend []uint32
	rcved uint32
}

type udpresult struct {
	serialNo uint64
	udpack
}

type UdpAcker interface {
	AppendResend(id...uint32)
	SetRcved(ids uint32)
	GetReSend() []uint32
	GetRcved() uint32
}

type UdpResulter interface {
	SetSerialNo(sn uint64)
	GetSerialNo() uint64
	UdpAcker
}

func NewUdpResult(sn uint64) UdpResulter {
	ur := &udpresult{serialNo:sn}
	ur.resend = make([]uint32,0)

	return ur
}

func (ur *udpresult) SetSerialNo(sn uint64) {
	ur.serialNo = sn
}

func (ur *udpresult) GetSerialNo() uint64 {
	return ur.serialNo
}

func (ua *udpack) AppendResend(ids...uint32) {
	for _,id:= range ids{
		ua.resend = append(ua.resend,id)
	}

}


func (ua *udpack) SetRcved(id uint32) {
	ua.rcved = id
}


func (ua *udpack) GetResend() []uint32 {
	return  ua.resend
}


func (ua *udpack) GetRcved() uint32 {
	return ua.rcved
}

