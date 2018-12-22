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

	if ipstr == "localaddress"{
		ua = address.GetAllLocalIPAddr(port)
	}else if ipstr == ""{
		ipstr = "0.0.0.0"
	}else {
		ua.AddIP4(ipstr,port)
	}
	fmt.Println("Server will start at:")
	ua.PrintAll()

	strip,p:= ua.FirstS()

	uaddr := &net.UDPAddr{IP:net.ParseIP(strip),Port:int(p)}

	socket,err := net.ListenUDP("udp4", uaddr)
	if err != nil {
		fmt.Println("监听失败!", err)
		return
	}
	defer socket.Close()

	for {
		// 读取数据
		data := make([]byte, 4096)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("读取数据失败!", err)
			continue
		}
		fmt.Println(read, remoteAddr)
		fmt.Printf("%s\n\n", data)

		// 发送数据
		senddata := []byte("hello client!")
		_, err = socket.WriteToUDP(senddata, remoteAddr)
		if err != nil {
			return
			fmt.Println("发送数据失败!", err)
		}
	}

	//listen

}

func (us *udpServer)Send([]byte) error {
	return nil
}

