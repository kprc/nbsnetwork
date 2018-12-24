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

	mListenSocket map[address.UdpAddresser]*net.UDPConn

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
	us.mListenSocket = make(map[address.UdpAddresser]*net.UDPConn)
	us.remoteAddr = make(map[address.UdpAddresser]*net.UDPAddr)
	//us.rcvBuf = make(map[address.UdpAddresser][]bytes.Buffer)

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


		//for {
		//
		//	data := make([]byte, 1024)
		//	read, remoteAddr, err := socket.ReadFromUDP(data)
		//	if err != nil {
		//		fmt.Println("读取数据失败!", err)
		//		continue
		//	}
		//	fmt.Println(read, remoteAddr)
		//	fmt.Printf("%s\n\n", data)
		//
		//	senddata := []byte("hello client!")
		//	_, err = socket.WriteToUDP(senddata, remoteAddr)
		//	if err != nil {
		//		fmt.Println("发送数据失败!", err)
		//		return
		//	}
		//}
	}

	us.listenAddr = usua

	usua.PrintAll()

	//listen

}

func (us *udpServer)Send([]byte) error {
	return nil
}

