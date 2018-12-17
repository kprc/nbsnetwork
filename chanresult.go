package nbsnetwork


type udpresult struct {
	serialNo uint64
	resend []uint32
	rcved []uint32
}


type UdpResulter interface {
	SetSerialNo(sn uint64)
	GetSerialNo() uint64
	AppendResend(ids...uint32)
	AppendRcved(ids...uint32)
}

func NewUdpResult(sn uint64) UdpResulter {
	ur := &udpresult{serialNo:sn}
	ur.resend = make([]uint32,0)
	ur.rcved = make([]uint32,0)

	return ur
}

func (ur *udpresult) SetSerialNo(sn uint64) {
	ur.serialNo = sn
}

func (ur *udpresult) GetSerialNo() uint64 {
	return ur.serialNo
}

func (ur *udpresult)AppendResend(ids...uint32)  {
	for _,id := range ids{
		ur.resend = append(ur.resend,id)
	}
}
func (ur *udpresult)AppendRcved(ids...uint32)  {
	for _,id := range ids{
		ur.rcved = append(ur.rcved,id)
	}
}