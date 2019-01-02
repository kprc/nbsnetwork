package client

import (
	"github.com/kprc/nbsnetwork/common/address"
	"io"
)

type udpClient struct {
	dialAddr address.UdpAddresser

	processWait chan int
}


type UdpClient interface {
	Send(r io.ReadSeeker) error
}

func NewUdpClient(ip string,port uint16) UdpClient {
	uc := &udpClient{dialAddr:address.NewUdpAddress(),processWait:make(chan int,1024)}
	uc.dialAddr.AddIP4(ip,port)

	return uc
}

func (uc *udpClient)Dial()  {
	
}

func (uc *udpClient)Send(r io.ReadSeeker) error  {



	return nil
}