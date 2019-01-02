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




