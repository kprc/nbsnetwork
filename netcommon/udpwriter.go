package netcommon

import (
	"github.com/kprc/nbsnetwork/common/constant"
	"io"
	"net"
	"github.com/kprc/nbsnetwork/send"
)

type udpReaderWriter struct {
	addr *net.UDPAddr
	sock *net.UDPConn
	needRemoteAddress bool
}


type UdpReaderWriterer interface {
	Send(r io.ReadSeeker) error
	SendBytes(data []byte)
	IsNeedRemoteAddress() bool
	NeedRemoteAddress()
	AddrString() string
	GetSock() *net.UDPConn
	GetAddr() *net.UDPAddr
	io.Writer
	io.Reader
}

func NewReaderWriter(addr *net.UDPAddr, sock *net.UDPConn) UdpReaderWriterer {
	uw:=&udpReaderWriter{addr:addr,sock:sock}

	return uw
}

func (uw *udpReaderWriter)GetSock() *net.UDPConn {
	return uw.sock
}

func (uw *udpReaderWriter)GetAddr() *net.UDPAddr  {
	return uw.addr
}

func (uw *udpReaderWriter)SendBytes(data []byte)  {
	uwrs := NewReadSeeker(data)

	uw.Send(uwrs)
}

func (uw *udpReaderWriter)AddrString() string  {
	return uw.addr.String()
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
	return uw.sock.Write(p)
}

func (uw *udpReaderWriter)Read(p []byte) (n int, err error)  {
	return uw.sock.Read(p)
}

func (uw *udpReaderWriter)IsNeedRemoteAddress() bool {
	return uw.needRemoteAddress
}

func (uw *udpReaderWriter) NeedRemoteAddress() {
	uw.needRemoteAddress = true
}