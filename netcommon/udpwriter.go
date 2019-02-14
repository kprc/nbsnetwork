package netcommon

import (
	"github.com/kprc/nbsdht/nbserr"
	"io"
	"net"
	"time"
)

type udpReaderWriter struct {
	addr *net.UDPAddr
	sock *net.UDPConn
	needRemoteAddress bool
	ok bool
}


type UdpReaderWriterer interface {
	IsNeedRemoteAddress() bool
	NeedRemoteAddress()
	SetNeedRemote(is bool)
	AddrString() string
	SetSock(sock *net.UDPConn)
	GetSock() *net.UDPConn
	SetAddr(addr *net.UDPAddr)
	GetAddr() *net.UDPAddr
	SetSockNull()
	SetDeadLine(tv time.Duration)
	IsTimeOut(err error) bool
	IsOK() bool
	ReadUdp(p []byte)(int,*net.UDPAddr,error)
	io.Writer
	io.Reader
	Close() error
}

func NewReaderWriter(addr *net.UDPAddr, sock *net.UDPConn,need bool) UdpReaderWriterer {
	uw:=&udpReaderWriter{addr:addr,sock:sock,needRemoteAddress:need}

	return uw
}

func (uw *udpReaderWriter)Close()  error{
	return uw.sock.Close()
}

func (uw *udpReaderWriter)SetSockNull()  {
	uw.sock = nil
}

func (uw *udpReaderWriter)SetNeedRemote(is bool)  {
	uw.needRemoteAddress = is
}

func (uw *udpReaderWriter)SetSock(sock *net.UDPConn)  {
	uw.sock = sock
}

func (uw *udpReaderWriter)GetSock() *net.UDPConn {
	return uw.sock
}

func (uw *udpReaderWriter)SetAddr(addr *net.UDPAddr)  {
	uw.addr = addr
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

func (uw *udpReaderWriter)ReadUdp(p []byte)(int,*net.UDPAddr,error)  {
	if uw.sock == nil {
		return 0,nil,nbserr.NbsErr{Errmsg:"sock is none",ErrId:nbserr.UDP_RCV_DEFAULT_ERR}
	}
	if uw.needRemoteAddress {
		return uw.sock.ReadFromUDP(p)
	}
	n,err:=uw.sock.Read(p)

	return n,uw.addr,err
}


func (uw *udpReaderWriter)IsNeedRemoteAddress() bool {
	return uw.needRemoteAddress
}

func (uw *udpReaderWriter) NeedRemoteAddress() {
	uw.needRemoteAddress = true
}

func (uw *udpReaderWriter)SetDeadLine(tv time.Duration)  {
	uw.sock.SetReadDeadline(time.Now().Add(time.Millisecond*tv))
}

func (uw *udpReaderWriter)IsTimeOut(err error) bool {
	if nerr,ok:=err.(net.Error);ok && nerr.Timeout(){
		uw.ok = false
		return true
	}
	return false
}

func (uw *udpReaderWriter)IsOK() bool {
	return uw.ok
}