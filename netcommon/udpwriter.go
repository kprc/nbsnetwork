package netcommon

import (
	"io"
	"net"
)

type udpReaderWriter struct {
	addr *net.UDPAddr
	sock *net.UDPConn
	needRemoteAddress bool
}


type UdpReaderWriterer interface {
	IsNeedRemoteAddress() bool
	NeedRemoteAddress()
	SetNeedRmAddr(is bool)
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

func (uw *udpReaderWriter)SetNeedRmAddr(is bool)  {
	uw.NeedRemoteAddress()
}

func (uw *udpReaderWriter)GetSock() *net.UDPConn {
	return uw.sock
}

func (uw *udpReaderWriter)GetAddr() *net.UDPAddr  {
	return uw.addr
}

func (uw *udpReaderWriter)AddrString() string  {
	return uw.addr.String()
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