package server

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/recv"
	"net"
	"sync"
)

var (
	instance UdpServerer
	instlock sync.Mutex
)


type udpServer struct {
	listenAddr address.UdpAddresser

	//mListenSocket map[address.UdpAddresser]*net.UDPConn

	//remoteAddr map[address.UdpAddresser]*net.UDPAddr

	processWait chan int

	//rcvBuf map[address.UdpAddresser][]bytes.Buffer
}

type UdpServerer interface {
	Run(ipstr string,port uint16)
	//Send([] byte) error
	GetListenAddr() address.UdpAddresser
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
	//us.mListenSocket = make(map[address.UdpAddresser]*net.UDPConn)
	//us.remoteAddr = make(map[address.UdpAddresser]*net.UDPAddr)
	//us.rcvBuf = make(map[address.UdpAddresser][]bytes.Buffer)
	us.processWait = make(chan int,0)



	return us
}


func (us *udpServer)Run(ipstr string,port uint16) {
	var ua address.UdpAddresser

	if ipstr == "localaddress" {
		ua = address.GetAllLocalIPAddr(port)
	} else {
		ua = address.NewUdpAddress()
		if ipstr == "" {
			ipstr = "0.0.0.0"
			ua.AddIP4(ipstr, port)
		} else {
			ua.AddIP4(ipstr, port)
		}
	}
	//fmt.Println("Server will start at:")
	//ua.PrintAll()

	usua:=address.NewUdpAddress()

	ua.Iterator()

	for {
		strip, p := ua.NextS()

		if strip == "" {
			break
		}


		netaddr := &net.UDPAddr{IP: net.ParseIP(strip), Port: int(p)}

		sock, err := net.ListenUDP("udp4", netaddr)
		if err != nil {
			fmt.Println("Listen Failure", err)
			return
		}

		//fmt.Println(netaddr.Port,netaddr.IP.String())
		//fmt.Println(s.LocalAddr().String())
		usua.AddIP4Str(sock.LocalAddr().String())


		go sockRecv(sock)
	}

	us.listenAddr = usua

	fmt.Println("Server will start at:")
	usua.PrintAll()

	wait:=<-us.processWait
	fmt.Println("receive a quit command",wait)

}

func (us *udpServer)GetListenAddr() address.UdpAddresser  {
	return us.listenAddr
}



func sockRecv(sock *net.UDPConn) error{
	dispatch:=recv.NewUdpDispath(true)
	dispatch.SetSock(sock)

	return dispatch.Dispatch()

}

