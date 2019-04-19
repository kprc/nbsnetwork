package netcommon

import (
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/tools"
	"net"
	"sync"
)


type udpconn struct {
	dialAddr address.UdpAddresser
	localAddr address.UdpAddresser
	realAddr address.UdpAddresser
	addr *net.UDPAddr
	sock *net.UDPConn
	isconn bool     //if conn from listen, isconn is false
	ready2send chan interface{}
	tick chan int64
	statuslock sync.Mutex
	status int32 //0 not set, 1 stopped, 2 bad connection, 3 connecting
	stopsendsign chan int
	closesign chan int
	lastrcvtime int64
	lastSendKaTime int64
	timeouttv int
	wg *sync.WaitGroup
	tickrcv chan int64
	rcv chan int
}

type UdpConn interface {
	Dial() error
	ConnSync()
	Connect() error
	Send(data [] byte,typ uint32) error
	SendAsync(data [] byte, typ uint32) error  //unblocking
	SetTimeout(tv int)
	Status() bool
	Close()
	GetAddr() address.UdpAddresser
	IsConnectTo(addr *net.UDPAddr) bool
	Push(v interface{})
	WaitConnReady() bool
	UpdateLastAccessTime()
}


//======= For Send Buffer Define Begin ======

type senddata struct {
	data []byte
	msgtyp uint32
}

//=======For Send Buffer Define End  =====

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
	uc.rcv = make(chan int)
	uc.tickrcv = make(chan int64)
	uc.isconn = true

	nt := tools.GetNbsTickerInstance()
	nt.Reg(&uc.tick)
	nt.RegWithTimeOut(&uc.tickrcv,1000)
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

	if uc.wg !=nil {
		uc.wg.Done()
	}

	for{
		select {
			case data2send:=<-uc.ready2send:
				d2s := data2send.(senddata)
				if err := uc.send(d2s.data,CONN_PACKET_TYP_DATA,d2s.msgtyp);err!=nil{

					uc.status = BAD_CONNECTION
					if uc.isconn == true {
						uc.sock.Close()
					}
					uc.sock = nil
					return err
				}

			case <-uc.tick:
				if err := uc.sendKAPacket();err!=nil{
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
				uc.status = STOP_CONNECTION
				closesign = 1
				return nil

		}
	}
}



func (uc *udpconn)recv(wg *sync.WaitGroup) error{
	defer wg.Done()

	n:=1024
	roundbuf := make([]byte,n)

	if uc.wg !=nil{
		uc.wg.Done()
	}

	for {
		buf := roundbuf[0:n]

		var err error
		var nr int
		if  nr,err = uc.sock.Read(buf); err!=nil{
			return baderr
		}

		uc.UpdateLastAccessTime()

		cp:=NewConnPacket()
		if err=cp.DeSerialize(buf[0:nr]);err!=nil{
			continue
		}

		if cp.GetTyp() == CONN_PACKET_TYP_KA{
			continue
		}

		uc.Push(cp)
	}

	return nil

}

func (uc *udpconn)UpdateLastAccessTime()  {

	if uc.isconn {
		if uc.lastrcvtime == 0 {
			uc.rcv <- 0
		}
	}

	uc.lastrcvtime = tools.GetNowMsTime()
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


	tv:=tools.GetNowMsTime() - uc.lastrcvtime

	if tv > int64(uc.timeouttv) && uc.lastrcvtime != 0{
		return baderr
	}

	tv=tools.GetNowMsTime() - uc.lastSendKaTime

	if  tv> int64(uc.timeouttv)/10 {
		return uc.send([]byte("ka message"), CONN_PACKET_TYP_KA,0)
	}


	return nil
}


func (uc *udpconn)send(v interface{}, typ uint32,msgtyp uint32) error {

	data:=v.([]byte)

	cp:=NewConnPacket()

	cp.SetTyp(typ)
	cp.SetMsgTyp(msgtyp)
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

	uc.lastSendKaTime = tools.GetNowMsTime()

	return nil
}

func (uc *udpconn)Send(data []byte,typ uint32) error  {
	if uc.status == CONNECTION_INIT || uc.status == STOP_CONNECTION {
		return notreadyerr
	}

	if uc.status == BAD_CONNECTION {
		return baderr
	}
	//copy?
	snddata:=senddata{data,typ}

	uc.ready2send <- snddata

	return nil
}

func (uc *udpconn)SendAsync(data []byte,typ uint32) error  {
	if uc.status == CONNECTION_INIT || uc.status == STOP_CONNECTION {
		return notreadyerr
	}

	if uc.status == BAD_CONNECTION {
		return baderr
	}

	snddata := senddata{data,typ}

	select {
	case uc.ready2send <- snddata:
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

func (uc *udpconn)ConnSync() {
	uc.wg = &sync.WaitGroup{}
	if uc.isconn {
		uc.wg.Add(2)
	} else {
		uc.wg.Add(1)
	}

}

func (uc *udpconn)WaitConnReady() bool {
	uc.wg.Wait()
	uc.wg = nil

	if uc.isconn {
		defer func() {
			nt:=tools.GetNbsTickerInstance()
			nt.UnReg(&uc.tickrcv)
		}()
		select {
		case <-uc.rcv:
			return true
		case <-uc.tickrcv:
			return false
		}
	}

	return true
}

