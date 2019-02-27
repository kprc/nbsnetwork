package netcommon

import (
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/tools"
	"net"
	"sync"
	"time"
)


type udpconn struct {
	dialAddr address.UdpAddresser
	localAddr address.UdpAddresser
	realAddr address.UdpAddresser
	addr *net.UDPAddr
	sock *net.UDPConn
	isconn bool     //if conn from listen, isconn is false
	ready2send chan interface{}
	//recvFromConn chan interface{}
	tick chan int64
	statuslock sync.Mutex
	status int32 //0 not set, 1 stopped, 2 bad connection, 3 connecting
	stopsendsign chan int
	closesign chan int
	lastrcvtime int64
	timeouttv int
	wg *sync.WaitGroup
}

type UdpConn interface {
	Dial() error
	Hello()
	Connect() error
	Send(data [] byte) error
	//Read() ([]byte, error)
	SendAsync(data [] byte) error  //unblocking
	//ReadAsync() ([]byte,error)     //unblocking
	SetTimeout(tv int)
	Status() bool
	Close()
	GetAddr() address.UdpAddresser
	IsConnectTo(addr *net.UDPAddr) bool
	Push(v interface{})
	WaitHello()
}


var (
	CONNECTION_INIT int32 = 0
	STOP_CONNECTION int32 = 1
	BAD_CONNECTION int32 = 2
	CONNECTION_RUNNING int32 = 3
)

var (
	dialerr  = nbserr.NbsErr{ErrId:nbserr.UDP_DIAL_ERR,Errmsg:"Dial UDP Connection Error"}
	baderr = nbserr.NbsErr{ErrId:nbserr.UDP_BAD_CONN,Errmsg:"Bad Connection"}
	notreadyerr = nbserr.NbsErr{ErrId:nbserr.UDP_CONN_NOTREADY,Errmsg:"Connection not ready!"}
	nodataerr = nbserr.NbsErr{ErrId:nbserr.UDP_CONN_NODATA,Errmsg:"Connection have no data arrived"}
	bufferoverflowerr = nbserr.NbsErr{ErrId:nbserr.UDP_BUFFOVERFLOW,Errmsg:"Buffer overflow"}
	listenconnerr = nbserr.NbsErr{ErrId:nbserr.UDP_CONN_LISTEN,Errmsg:"Connection is received"}
)

func NewUdpConnFromListen(addr *net.UDPAddr,sock *net.UDPConn) UdpConn {
	uc:= &udpconn{}

	uc.ready2send = make(chan interface{},1024)
	//uc.recvFromConn = make(chan interface{},1024)
	uc.tick = make(chan int64,8)
	uc.stopsendsign = make(chan int)
	uc.closesign = make(chan int)

	uc.addr = addr
	uc.sock = sock
	uc.isconn = false

	nt := tools.GetNbsTickerInstance()
	nt.Reg(&uc.tick)
	uc.timeouttv = 10000   //ms



	return uc

}

func NewUdpCreateConnection(rip,lip string,rport,lport uint16) UdpConn  {
	if rip == ""{
		return nil
	}
	uc:= &udpconn{}
	uc.dialAddr = address.NewUdpAddress()
	uc.dialAddr.AddIP4(rip,rport)
	if lip != "" {
		uc.localAddr = address.NewUdpAddress()
		uc.localAddr.AddIP4(lip,lport)
	}else if lport != 0{
		uc.localAddr.AddIP4("0.0.0.0",lport)
	}

	uc.ready2send = make(chan interface{},1024)
	//uc.recvFromConn = make(chan interface{},1024)
	uc.tick = make(chan int64,8)
	uc.stopsendsign = make(chan int)
	uc.closesign = make(chan int)
	uc.isconn = true

	nt := tools.GetNbsTickerInstance()
	nt.Reg(&uc.tick)
	uc.timeouttv = 10000   //ms


	return uc
}


func assembleUdpAddr(addr address.UdpAddresser) *net.UDPAddr  {
	if addr == nil {
		return nil
	}
	ipstr,port := addr.FirstS()
	if ipstr == "" {
		return nil
	}
	return &net.UDPAddr{IP:net.ParseIP(ipstr),Port:int(port)}
}

func (uc *udpconn)Dial() error  {
	if uc.isconn == false{
		return listenconnerr
	}
	la := assembleUdpAddr(uc.localAddr)
	ra := assembleUdpAddr(uc.dialAddr)

	conn,err:=net.DialUDP("udp4",la,ra)
	if err!=nil{
		return dialerr
	}

	if la != nil{
		uc.realAddr = address.NewUdpAddress()
		uc.realAddr.AddIP4Str(la.String())
	}

	uc.sock = conn
	uc.addr = ra

	return nil


}

func (uc *udpconn)SetTimeout(tv int)  {
	uc.timeouttv = tv
}

func (uc *udpconn)Connect() error{

	var closesign int = 0

	uc.statuslock.Lock()


	if uc.status == BAD_CONNECTION{
		uc.statuslock.Unlock()
		return baderr
	}

	if uc.status == CONNECTION_RUNNING {
		uc.statuslock.Unlock()
		return nil
	}

	uc.status = CONNECTION_RUNNING

	uc.statuslock.Unlock()
	if uc.wg !=nil {
		uc.wg.Done()
	}

	if uc.isconn {
		wg := &sync.WaitGroup{}
		//start send and recv
		wg.Add(1)

		go uc.recv(wg)

		defer func() {
			//fmt.Println("begin defer func")
			wg.Wait()
			//fmt.Println("close sign",closesign)
			if closesign == 1 {
				//fmt.Println("send close sign .")
				uc.closesign <- 1
			}

		}()
	}

	for{
		select {
			case data2send:=<-uc.ready2send:
				if err := uc.send(data2send,CONN_PACKET_TYP_DATA);err!=nil{

					uc.status = BAD_CONNECTION
					if uc.isconn == true {
						uc.sock.Close()
					}
					uc.sock = nil
					return err
				}

			case <-uc.tick:
				//fmt.Println("get tick",t)
				if err := uc.sendKAPacket();err!=nil{
					//fmt.Println("get ka send err")
					uc.status = BAD_CONNECTION
					if uc.isconn == true {
						if uc.sock != nil{
							uc.sock.Close()
						}
					}
					uc.sock = nil
					return err
				}
			case <-uc.stopsendsign:
				//fmt.Println("Receive stop sign in connect func")
				uc.status = STOP_CONNECTION
				closesign = 1
				return nil

		}
	}
}

func getNowMsTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func (uc *udpconn)recv(wg *sync.WaitGroup) error{
	defer wg.Done()

	//fmt.Println("start recv")

	n:=1024
	roundbuf := make([]byte,n)

	for {
		buf := roundbuf[0:n]

		var err error
		var nr int
		if  nr,err = uc.sock.Read(buf); err!=nil{
			//uc.sock.Close()
			//uc.sock = nil
			//fmt.Println("recv sock read err")
			return baderr
		}

		uc.lastrcvtime =  getNowMsTime()

		cp:=NewConnPacket()
		if err=cp.DeSerialize(buf[0:nr]);err!=nil{
			continue
		}

		if cp.GetTyp() == CONN_PACKET_TYP_KA{
			continue
		}

		uc.Push(cp)
	}
	//fmt.Println("recv quit")

	return nil

}


func (uc *udpconn)Close() {

	if uc.status == BAD_CONNECTION || uc.status == STOP_CONNECTION{
		return
	}
	if uc.status == CONNECTION_RUNNING {
		uc.stopsendsign <- 0
		if uc.sock != nil {
			if uc.isconn == true {
				uc.sock.Close()
			}
		}
		//fmt.Println("send stop send sign done")
		if uc.isconn == true {
			<-uc.closesign
		}

	}else {
		if uc.sock != nil {
			if uc.isconn == true {
				uc.sock.Close()
			}
		}

	}

	uc.sock = nil
}

func (uc *udpconn)sendKAPacket() error {
	//if uc.isconn == false{
	//	return nil
	//}
	tv:=getNowMsTime() - uc.lastrcvtime
	if  tv> int64(uc.timeouttv)/10 {
		if tv < int64(uc.timeouttv) || uc.lastrcvtime == 0{
			return uc.send([]byte("ka message"), CONN_PACKET_TYP_KA)
		}else{
			return baderr
		}
	}

	return nil
}


func (uc *udpconn)send(v interface{}, typ uint32) error {
	data:=v.([]byte)

	//fmt.Println("udp conn send msg:",string(data))

	cp:=NewConnPacket()

	cp.SetTyp(typ)
	cp.SetData(data)
	id:=nbsid.GetLocalId()
	cp.SetUid(id.Bytes())

	var d []byte
	var err error

	if d,err=cp.Serialize();err!=nil{
		return nil
	}

	//send d
	if uc.isconn {
		if _,err1:=uc.sock.Write(d);err1!=nil{
			return baderr
		}
	}else {
		if _,err1:=uc.sock.WriteToUDP(d,uc.addr);err1!=nil{
			return baderr
		}
	}

	return nil
}

func (uc *udpconn)Send(data []byte) error  {
	if uc.status == CONNECTION_INIT || uc.status == STOP_CONNECTION {
		return notreadyerr
	}

	if uc.status == BAD_CONNECTION {
		return baderr
	}
	//copy?
	uc.ready2send <- data

	return nil
}

func (uc *udpconn)SendAsync(data []byte) error  {
	if uc.status == CONNECTION_INIT || uc.status == STOP_CONNECTION {
		return notreadyerr
	}

	if uc.status == BAD_CONNECTION {
		return baderr
	}

	select {
	case uc.ready2send <- data:
		return nil
	default:
		return bufferoverflowerr
	}
}

func (uc *udpconn)Status() bool  {
	if uc.status == CONNECTION_RUNNING {
		return true
	}

	return false
}

//func (uc *udpconn)Read() ([]byte,error)  {
//
//	if uc.status == CONNECTION_INIT || uc.status == STOP_CONNECTION {
//		return nil,notreadyerr
//	}
//
//	if uc.status == BAD_CONNECTION {
//		return nil,baderr
//	}
//
//
//	ret := <-uc.recvFromConn
//
//	return ret.([]byte),nil
//}
//
//func (uc *udpconn)ReadAsync() ([]byte,error)  {
//	if uc.status == CONNECTION_INIT || uc.status == STOP_CONNECTION {
//		return nil,notreadyerr
//	}
//
//	if uc.status == BAD_CONNECTION {
//		return nil,baderr
//	}
//
//	select {
//	case ret:=<-uc.recvFromConn:
//		return ret.([]byte),nil
//	default:
//		return nil,nodataerr
//	}
//}

func (uc *udpconn)GetAddr() address.UdpAddresser  {
	addr:=address.NewUdpAddress()

	addr.AddIP4Str(uc.addr.String())

	return addr
}

func (uc *udpconn)IsConnectTo(addr *net.UDPAddr) bool {
	ipstr := uc.addr.String()

	ipstr2 := addr.String()


	if ipstr == ipstr2{
		return true
	}

	return false

}

func (uc *udpconn)Push(v interface{})  {
	cs:=GetConnStoreInstance()

	rb:=NewRcvBlock(v.(ConnPacket),uc)

	cs.Push(rb)
}

func (uc *udpconn)Hello()  {
	uc.wg = &sync.WaitGroup{}
	uc.wg.Add(1)
}

func (uc *udpconn)WaitHello()  {
	uc.wg.Wait()
	uc.wg = nil
}
