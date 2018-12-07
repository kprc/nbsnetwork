package nbsnetwork

import "sync/atomic"

type udpHead struct {
	serialNo uint64
	totalCnt uint32
	posNum uint32
	data interface{}
}


type UdpHeader interface {
	nextSerialNo() uint64
}



var gSerialNo uint64 = UDP_SERIAL_MAGIC_NUM

func (uh *udpHead)nextSerialNo() uint64 {
	return atomic.AddUint64(&gSerialNo,1)
}

func NewUdpHead()  UdpHeader {
	uh := udpHead{}
	uh.nextSerialNo()

	return &uh
}

func (uh *udpHead)Byte() []byte  {
	return nil
}

func (uh *udpHead)SetPos(pos uint32)  {
	atomic.StoreUint32(&uh.posNum,pos)
}

func (uh *udpHead)IncPos() uint32 {
	return atomic.AddUint32(&uh.posNum,1)
}

func (uh *udpHead)GetPos() uint32  {
	return atomic.LoadUint32(&uh.posNum)
}