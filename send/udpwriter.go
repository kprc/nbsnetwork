package send

import (
	"github.com/kprc/nbsnetwork/common/constant"
	"io"
	"net"
)

type udpWriter struct {
	addr *net.UDPAddr
	sock *net.UDPConn
}


type UdpWriterer interface {
	Send(r io.ReadSeeker) error
	io.Writer
}

func NewWriter(addr *net.UDPAddr, sock *net.UDPConn) UdpWriterer {
	uw:=&udpWriter{addr:addr,sock:sock}

	return uw
}

func (uw *udpWriter)Send(r io.ReadSeeker) error  {
	bd := NewBlockData(r,constant.UDP_MTU)
	bd.SetWriter(uw)

	sd:=NewStoreData(bd)

	bs:=GetInstance()

	bs.AddBlockDataer(bd.GetSerialNo(),sd)

	sd = bs.GetBlockDataer(bd.GetSerialNo())

	sd.GetBlockData().Send()

	return nil
}

func (uw *udpWriter)Write(p []byte) (n int, err error)   {
	return uw.sock.WriteToUDP(p,uw.addr)
}