package netcommon

import (
	"github.com/kprc/nbsdht/nbserr"
	"net"
	"sync"
)

type udpconn struct {
	addr *net.UDPAddr
	sock *net.UDPConn
	isconn bool     //if conn from listen, isconn is false
	ready2send chan interface{}
	recvFromConn chan interface{}

	status int32 //0 not set, 1 stopped, 2 bad connection,

	runninglock sync.Mutex
	isrunning int32 //0 not running, 1 running

}

 var (
 	baderr = nbserr.NbsErr{ErrId:nbserr.UDP_BAD_CONN,Errmsg:"Bad Connection"}
 )


func (uc *udpconn)Connect() error {
	if uc.status == 2{
		return baderr
	}

	//start send and recv



}


func (uc *udpconn)Send(v interface{}) error  {
	return nil
}



