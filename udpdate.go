package nbsnetwork

import "sync/atomic"

type udpHead struct {
	serialNo uint64
	totalCnt uint
	posNum uint
	data interface{}
}


type UdpHeader interface {
	nextSerialNo() uint64
}



var gSerialNo uint64 = UDP_SERIAL_MAGIC_NUM

func (uh *udpHead)nextSerialNo() uint64 {
	uh.serialNo = atomic.AddUint64(&gSerialNo,1)
}

func NewUdpHead()  UdpHeader {
	uh := udpHead{}
	uh.nextSerialNo()

	return &uh
}

func (uh *udpHead)Byte() []byte  {
	return nil
}

