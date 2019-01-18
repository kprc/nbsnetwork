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
	SetNeedRemote(is bool)
	AddrString() string
	GetSock() *net.UDPConn
	GetAddr() *net.UDPAddr
	SetSockNull()
	io.Writer
	io.Reader
}

func NewReaderWriter(addr *net.UDPAddr, sock *net.UDPConn,need bool) UdpReaderWriterer {
	uw:=&udpReaderWriter{addr:addr,sock:sock,needRemoteAddress:need}

	return uw
}

func (uw *udpReaderWriter)SetSockNull()  {
	uw.sock = nil
}

func (uw *udpReaderWriter)SetNeedRemote(is bool)  {
	uw.needRemoteAddress = is
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

	if uw.IsNeedRemoteAddress(){
		return uw.sock.WriteToUDP(p, uw.addr)
	}else {
		return uw.sock.Write(p)
	}
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