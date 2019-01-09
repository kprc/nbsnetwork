package server

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/flowkey"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/recv"
	"github.com/kprc/nbsnetwork/send"
	"github.com/kprc/nbsnetwork/common/message"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"io"
	"net"
	"sync"
	"github.com/kprc/nbsnetwork/rw"
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

func doAck(pkt packet.UdpPacketDataer)  {
	bs:= send.GetInstance()

	rmr := regcenter.GetMsgCenterInstance()
	_,stationId,_:=rmr.GetMsgId(pkt.GetTransInfo())

	if stationId == "" {
		return
	}

	fk := flowkey.NewFlowKey(stationId,pkt.GetSerialNo())

	sd := bs.GetBlockDataer(fk)
	if sd == nil {
		return
	}

	bdr := sd.GetBlockData()
	bdr.PushResult(pkt.GetData())
	bs.PutBlockDataer(fk)

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

		//if type is ack, send to result channel and continue
		if pkt.GetTyp() == constant.ACK {
			//send to send block
			doAck(pkt)
			continue
		}

		mc := regcenter.GetMsgCenterInstance()
		//get transfer layer information
		msgid,stationId,headinfo := mc.GetMsgId(pkt.GetTransInfo())
		if msgid == constant.MSG_NONE || stationId == "" {
			fmt.Printf("Receive msgid %d, stationId %s",msgid,stationId)
			continue
		}
		//get msg block
		rmr := message.GetInstance()
		k := flowkey.NewFlowKey(stationId,pkt.GetSerialNo())
		m := rmr.GetMsg(k)
		var rcv recv.RcvDataer
		if m==nil {
			m = message.NewRcvMsg()
			handler := mc.GetHandler(msgid)
			m.SetWS(handler.GetWSNew()(headinfo))
			uw:=rw.NewReaderWriter(remoteAddr,sock)
			m.SetUW(uw)
			m.SetKey(k)
			rcv = recv.NewRcvDataer(pkt.GetSerialNo(),m.GetWS())
			m.SetRecv(rcv)

			rmr.AddMSG(k,m)

			m = rmr.GetMsg(k)   //just to increase the ref count
		}else {
			rcv = m.GetRecv()
		}
		//wait to test
		ack,err := rcv.Write(pkt)
		if ack != nil{
			//wait to test
			sendAck(m.GetUW(),ack,pkt)
		}

		if rcv.CanDelete() {
			rmr.PutMsg(k)
			rmr.DelMsg(k)
		}else {
			if rcv.Finish() {
				mc.GetHandler(msgid).GetHandler()(headinfo, m.GetWS(), m.GetUW())
				mc.PutHandler(msgid)
				rmr.PutMsg(k)
			} else {
				rmr.PutMsg(k)
			}
		}

	}
}

func sendAck(w io.Writer, ack packet.UdpResulter, pkt packet.UdpPacketDataer){
	upd := packet.NewUdpPacketData()
	upd.SetSerialNo(pkt.GetSerialNo())
	upd.SetTransInfo(pkt.GetTransInfo())
	upd.SetACK()
	back,_ := ack.Serialize()
	upd.SetLength(int32(len(back)))
	upd.SetData(back)

	w.Write(back)

}

