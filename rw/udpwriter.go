package rw

import (
	"github.com/kprc/nbsnetwork/common/constant"
	"io"
	"net"
	"github.com/kprc/nbsnetwork/send"
)

type udpReaderWriter struct {
	addr *net.UDPAddr
	sock *net.UDPConn
}


type uwReaderSeeker struct {
	data []byte
}

type UdpBytesReaderSeeker interface {
	io.ReadSeeker
}

func (uwrs *uwReaderSeeker)Read(p []byte) (n int, err error){
	minlen := len(p)
	if minlen > len(uwrs.data) {
		minlen = len(uwrs.data)
	}
	cplen := copy(p[0:minlen],uwrs.data)
	uwrs.data = uwrs.data[cplen:]
	return cplen,nil
}

func (uwrs *uwReaderSeeker)Seek(offset int64, whence int) (int64, error)  {
	return 0,nil
}

func NewReadSeeker(data []byte) UdpBytesReaderSeeker {
	return &uwReaderSeeker{data:data}
}


type UdpReaderWriterer interface {
	Send(r io.ReadSeeker) error
	SendBytes(data []byte)
	io.Writer
	io.Reader
}

func NewReaderWriter(addr *net.UDPAddr, sock *net.UDPConn) UdpReaderWriterer {
	uw:=&udpReaderWriter{addr:addr,sock:sock}

	return uw
}

func (uw *udpReaderWriter)SendBytes(data []byte)  {
	uwrs := NewReadSeeker(data)

	uw.Send(uwrs)
}

func (uw *udpReaderWriter)Send(r io.ReadSeeker) error  {
	bd := send.NewBlockData(r,constant.UDP_MTU)
	bd.SetWriter(uw)

	sd:=send.NewStoreData(bd)

	bs:=send.GetInstance()

	bs.AddBlockDataer(bd.GetSerialNo(),sd)

	sd = bs.GetBlockDataer(bd.GetSerialNo())

	sd.GetBlockData().Send()

	return nil
}

func (uw *udpReaderWriter)Write(p []byte) (n int, err error)   {
	return uw.sock.WriteToUDP(p,uw.addr)
}

func (uw *udpReaderWriter)Read(p []byte) (n int, err error)  {
	return uw.sock.Read(p)
}