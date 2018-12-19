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

func (us *udpServer)Run(ipstr string,port uint16)  {

	if ipstr == ""{

	}

	if port == 0 {

	}

}

