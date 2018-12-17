package nbsnetwork

import (
	"bytes"

	"net"

)




type udpServer struct {
	listenAddr UdpAddresser

	mconn map[UdpAddresser]*net.UDPConn

	remoteAddr map[UdpAddresser]*net.UDPAddr

	rcvBuf map[UdpAddresser][]bytes.Buffer
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

	var xx *net.UDPConn
	xx.Write()

	//c ,err := net.ListenUDP("udp4",&net.UDPAddr{
	//})
}
