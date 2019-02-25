package netcommon

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/tools"
	"net"
	"sync"
	)

var (
	CONNECTION_INIT int32 = 0
	STOP_CONNECTION int32 = 1
	BAD_CONNECTION int32 = 2
	CONNECTION_RUNNING int32 = 3
)

type udpconn struct {
	addr *net.UDPAddr
	sock *net.UDPConn
	isconn bool     //if conn from listen, isconn is false
	ready2send chan interface{}
	recvFromConn chan interface{}
	tick chan int64

	statuslock sync.Mutex
	status int32 //0 not set, 1 stopped, 2 bad connection, 3 connecting

	stopsign chan int

	lastrcvtime int64

	timeouttv int

}

type UdpConn interface {
	Connect() error
	Send(data [] byte) error
	Read() ([]byte, error)
	SetTimeout(tv int)
	Status() bool
}

var (
	baderr = nbserr.NbsErr{ErrId:nbserr.UDP_BAD_CONN,Errmsg:"Bad Connection"}
	notreadyerr = nbserr.NbsErr{ErrId:nbserr.UDP_CONN_NOTREADY,Errmsg:"Connection not ready!"}
	nodataerr = nbserr.NbsErr{ErrId:nbserr.UDP_CONN_NODATA,Errmsg:"Connection have no data arrived"}
)

func NewUdpConn(addr *net.UDPAddr,sock *net.UDPConn,isconn bool) UdpConn {
	uc:= &udpconn{}

	uc.ready2send = make(chan interface{},1024)
	uc.recvFromConn = make(chan interface{},1024)
	uc.tick = make(chan int64,8)
	uc.stopsign = make(chan int)

	uc.addr = addr
	uc.sock = sock
	uc.isconn = isconn

	nt := tools.GetNbsTickerInstance()
	nt.Reg(&uc.tick)
	uc.timeouttv = 1000   //ms

	return uc

}

func (uc *udpconn)SetTimeout(tv int)  {
	uc.timeouttv = tv
}

func (uc *udpconn)Connect() error{
	if uc.status == 2{
		return baderr
	}

	if uc.isconn == false {
		wg := &sync.WaitGroup{}
		//start send and recv
		wg.Add(1)

		go uc.recv(wg)

		defer wg.Wait()
	}

	for{
		select {
			case data2send:=<-uc.ready2send:
				if err := uc.sendpacket(data2send);err!=nil{
					return err
				}

			case <-uc.tick:
				if err := uc.sendkapacket();err!=nil{
					return err
				}
			case <-uc.stopsign:
				uc.stop()
				return nil

		}
	}
}


func (uc *udpconn)recv(wg *sync.WaitGroup) error{
	defer wg.Done()

	n:=1024
	roundbuf := make([]byte,n)

	for {
		buf := roundbuf[0:n]

		var err error
		var nr int
		if  nr,err = uc.sock.Read(buf); err!=nil{
			return err
		}
		uc.recvFromConn <- buf[0:nr]
	}

	return nil

}


func (uc *udpconn)sendkapacket() error {

}

func (uc *udpconn)stop() {

}

func (uc *udpconn)sendpacket(v interface{})  error{

}


func (uc *udpconn)Send(data []byte) error  {

	return nil
}

func (uc *udpconn)Status() bool  {

}

func (uc *udpconn)Read() ([]byte,error)  {

	if uc.status != CONNECTION_RUNNING {
		return nil,notreadyerr
	}

	ret := <-uc.recvFromConn

	return ret.([]byte),nil
}

func (uc *udpconn)ReadyAnSyc() ([]byte,error)  {
	if uc.status != CONNECTION_RUNNING {
		return nil,notreadyerr
	}

	select {
	case ret:=<-uc.recvFromConn:
		return ret.([]byte),nil
	default:
		return nil,nodataerr
	}
}

