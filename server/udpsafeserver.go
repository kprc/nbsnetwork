package server

import (
	"net"
	"github.com/kprc/nbsnetwork/common/address"
	"sync"
	"fmt"
)

var (
	instance UdpServerer
	instlock sync.Mutex
)


type udpServer struct {
	listenAddr address.UdpAddresser

	mconn map[address.UdpAddresser]*net.UDPConn

	remoteAddr map[address.UdpAddresser]*net.UDPAddr

	//rcvBuf map[address.UdpAddresser][]bytes.Buffer
}

type UdpServerer interface {
	Run(ipstr string,port uint16)
	Send([] byte) error
	//Rcv() ([]byte,error)
}

func GetUdpServer() UdpServerer {
	if instance == nil{
		instlock.Lock()
		if instance ==nil {
			instance = newUdpServer()
		}
		instlock.Unlock()
	}

	return instance
}

func newUdpServer() UdpServerer {
	us := &udpServer{}
	us.mconn = make(map[address.UdpAddresser]*net.UDPConn)
	us.remoteAddr = make(map[address.UdpAddresser]*net.UDPAddr)
	//us.rcvBuf = make(map[address.UdpAddresser][]bytes.Buffer)

	return us
}

func (us *udpServer)Run(ipstr string,port uint16)  {

	var ua address.UdpAddresser

	if ipstr == ""{
		ua = address.GetAllLocalIPAddr(port)
	}else {
		ua.AddIP4(ipstr,port)
	}
	fmt.Println("Server will start at:")
	ua.PrintAll()


	//listen

}

func (us *udpServer)Send([]byte) error {
	return nil
}

