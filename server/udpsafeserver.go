package server

import (
	"net"
	"github.com/kprc/nbsnetwork/common/address"
	"sync"
	"fmt"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/server/regcenter"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/server/message"
	"github.com/kprc/nbsnetwork/recv"
)

var (
	instance UdpServerer
	instlock sync.Mutex
)


type udpServer struct {
	listenAddr address.UdpAddresser

	mListenSocket map[address.UdpAddresser]*net.UDPConn

	remoteAddr map[address.UdpAddresser]*net.UDPAddr

	processWait chan int

	//rcvBuf map[address.UdpAddresser][]bytes.Buffer
}

type UdpServerer interface {
	Run(ipstr string,port uint16)
	Send([] byte) error
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
	us.mListenSocket = make(map[address.UdpAddresser]*net.UDPConn)
	us.remoteAddr = make(map[address.UdpAddresser]*net.UDPAddr)
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

	usua.PrintAll()

	fmt.Println("Server will start at:")

	wait:=<-us.processWait
	fmt.Println("receive a quit command",wait)

}

func (us *udpServer)GetListenAddr() address.UdpAddresser  {
	return us.listenAddr
}

func sockRecv(sock *net.UDPConn){
	gdata := make([]byte, 1024)

	for {

		data := gdata[0:]
		n, remoteAddr, err := sock.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Read data error", err)
			continue
		}

		pkt:=packet.NewUdpPacketData()
		if err = pkt.DeSerialize(data[:n]);err!=nil{
			fmt.Println("Packet DeSerialize failed!",err)
			continue
		}

		mc := regcenter.GetMsgCenter()

		msgid,stationId := mc.GetMsgId(pkt.GetTransInfo())
		if msgid == constant.MSG_NONE || stationId == "" {
			fmt.Printf("Receive msgid %d, stationId %s",msgid,stationId)
			continue
		}

		rmr := message.GetInstance()
		k := message.NewMsgKey(pkt.GetSerialNo(),stationId)
		m := rmr.GetMsg(k)
		var rcv recv.RcvDataer
		if m==nil {
			m = message.NewRcvMsg()
			handler := mc.GetHandler(msgid)
			m.SetWS(handler.GetWSNew()())
			m.SetAddr(remoteAddr)
			m.SetKey(k)
			m.SetSock(sock)
			rcv = recv.NewRcvDataer(pkt.GetSerialNo(),m.GetWS())
			m.SetRecv(rcv)

			rmr.AddMSG(k,m)

			m = rmr.GetMsg(k)   //just to increate the ref count
		}else {
			rcv = m.GetRecv()
		}

		ack,err := rcv.Write(pkt)
		if ack != nil{
			back,_ := ack.Serialize()
			m.GetSock().WriteToUDP(back,m.GetAddr())
		}
		if rcv.Finish() {
			mc.GetHandler(msgid).GetHandler()(m.GetWS(),m.GetSock())
		}

	}
}

func (us *udpServer)Send([]byte) error {
	return nil
}

