package server

import (
	"bytes"

	"net"

	"github.com/kprc/nbsnetwork/common/address"
)




type udpServer struct {
	listenAddr address.UdpAddresser

	mconn map[address.UdpAddresser]*net.UDPConn

	remoteAddr map[address.UdpAddresser]*net.UDPAddr

	rcvBuf map[address.UdpAddresser][]bytes.Buffer
}

type UdpServerer interface {
	Run()
	Send([] byte) error
	//Rcv() ([]byte,error)
}

func (us *udpServer)Run()  {

	//var x *net.TCPConn
	//
	//x.Write()

	//var xx *net.UDPConn
	//xx.Write()

	//c ,err := net.ListenUDP("udp4",&net.UDPAddr{
	//})
}
