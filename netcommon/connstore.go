package netcommon

import (
	"net"
	"sync"
)

type connblock struct {
	conn UdpConn
}

type connstore struct {
	store map[string]*connblock
	rcvpacket chan interface{}
}


type ConnStore interface {
	Add(uid string,conn UdpConn)
	Update(uid string, sock *net.UDPConn,addr *net.UDPAddr)
	Del(uid string)
	GetConn(uid string) UdpConn
	Push(v interface{}) error
	Read() RcvBlock
	ReadAsync() (RcvBlock,error)
}


var (
	storeinstance ConnStore
	glock sync.Mutex

)


func newConnStore() ConnStore {
	cs := &connstore{}
	cs.store = make(map[string]*connblock)
	cs.rcvpacket = make(chan interface{},1024)

	return cs
}

func GetConnStoreInstance() ConnStore  {
	if storeinstance != nil{
		return storeinstance
	}

	glock.Lock()
	defer glock.Unlock()
	if storeinstance != nil{
		return storeinstance
	}

	storeinstance = newConnStore()

	return storeinstance

}


func (cs *connstore)Add(uid string,conn UdpConn)  {
	v,ok:=cs.store[uid]
	if !ok{
		cs.store[uid] = &connblock{conn:conn}
		return
	}

	if !v.conn.Status(){
		cs.store[uid] = &connblock{conn:conn}
		return
	}

	addr1 := v.conn.GetAddr()
	addr2 := conn.GetAddr()

	if addr1.Equals(addr2){
		return
	}

	v.conn.Close()

	cs.store[uid] = &connblock{conn:conn}
}

func NewConn(sock *net.UDPConn,addr *net.UDPAddr) UdpConn {
	uc:=NewUdpConnFromListen(addr,sock)
	uc.ConnSync()
	go uc.Connect()
	uc.WaitConnReady()

	return uc
}

func (cs *connstore)Update(uid string, sock *net.UDPConn,addr *net.UDPAddr)  {
	v,ok:=cs.store[uid]
	if !ok{
		cs.store[uid] = &connblock{conn:NewConn(sock,addr)}
		return
	}

	if !v.conn.Status(){
		cs.store[uid] = &connblock{conn:NewConn(sock,addr)}
		return
	}

	if v.conn.IsConnectTo(addr){
		v.conn.UpdateLastAccessTime()
		return
	}

	v.conn.Close()


	cs.store[uid] = &connblock{conn:NewConn(sock,addr)}


}


func (cs *connstore)Del(uid string)  {
	v,ok:=cs.store[uid]
	if !ok{
		return
	}
	v.conn.Close()
	delete(cs.store,uid)
}

func (cs *connstore)GetConn(uid string) UdpConn {

	if v,ok:=cs.store[uid];ok{

		return v.conn
	}

	return nil
}



func (cs *connstore)Push(v interface{}) error{
	select {
		case cs.rcvpacket <-v:
			return nil
		default:
			return bufferoverflowerr
	}
}


func (cs *connstore)Read() RcvBlock{
	cp:=<-cs.rcvpacket

	return cp.(RcvBlock)
}

func (cs *connstore)ReadAsync() (RcvBlock,error){
	select {
		case cp:=<-cs.rcvpacket:
			return cp.(RcvBlock),nil
	    default:
			return nil,nodataerr

	}
}
