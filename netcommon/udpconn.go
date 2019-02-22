package netcommon

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/tools"
	"net"
	"sync"
	)

type udpconn struct {
	addr *net.UDPAddr
	sock *net.UDPConn
	isconn bool     //if conn from listen, isconn is false
	ready2send chan interface{}
	recvFromConn chan interface{}
	tick chan int64

	status int32 //0 not set, 1 stopped, 2 bad connection,

	runninglock sync.Mutex
	isrunning int32 //0 not running, 1 running

	stopsign chan int

	lastrcvtime int64

	timeouttv int

}

type UdpConn interface {
	Connect() error
	Send(v interface{}) error
	SetTimeout(tv int)
}

 var (
 	baderr = nbserr.NbsErr{ErrId:nbserr.UDP_BAD_CONN,Errmsg:"Bad Connection"}

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

	nt:=tools.GetNbsTickerInstance()
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

	//start send and recv
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

func (uc *udpconn)sendkapacket() error {

}

func (uc *udpconn)stop() {

}

func (uc *udpconn)sendpacket(v interface{})  error{

}


func (uc *udpconn)Send(v interface{}) error  {

	return nil
}



