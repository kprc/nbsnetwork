package netcommon

import (
	"net"
	"sync"
)

type connblock struct {
	conn UdpConn
	//configure
}

type connstore struct {
	store map[string]*connblock
}


type ConnStore interface {
	Add(uid string,conn UdpConn)
	Update(uid string, sock *net.UDPConn,addr *net.UDPAddr)
	Del(uid string)
	GetConn(uid string) UdpConn
}


var (
	storeinstance ConnStore
	glock sync.Mutex
)


func newConnStore() ConnStore {
	cs := &connstore{}
	cs.store = make(map[string]*connblock)

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
	uc.Connect()

	return uc
}

func (cs *connstore)Update(uid string, sock *net.UDPConn,addr *net.UDPAddr)  {
	v,ok:=cs.store[uid]
	if ok{
		cs.store[uid] = &connblock{conn:NewConn(sock,addr)}
		return
	}

	if !v.conn.Status(){
		cs.store[uid] = &connblock{conn:NewConn(sock,addr)}
		return
	}

	if v.conn.IsConnectTo(addr){
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

