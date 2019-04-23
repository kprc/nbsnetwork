package netcommon

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/address"
	"net"
	"sync"
)

type udplisten struct {
	listenAddr address.UdpAddresser
	sockarr []*net.UDPConn
	wg *sync.WaitGroup
}


type UpdListen interface {
	Run(ipstr string,port uint16)
	GetListenAddr() address.UdpAddresser
	Close()
}


var (
	udplisteninstance UpdListen
	glistenlock sync.Mutex
)

func newUdpListen() UpdListen {
	return &udplisten{wg:&sync.WaitGroup{},sockarr:make([]*net.UDPConn,0)}
}

func GetUpdListenInstance() UpdListen  {
	if udplisteninstance != nil{
		return udplisteninstance
	}

	glistenlock.Lock()
	defer glistenlock.Unlock()

	if udplisteninstance != nil{
		return udplisteninstance
	}

	udplisteninstance = newUdpListen()

	return udplisteninstance

}


func (ul *udplisten)Run(ipstr string,port uint16)  {
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

	usua:=address.NewUdpAddress()

	ua.Iterator()

	for {
		strip, p,_ := ua.NextS()

		if strip == "" {
			break
		}


		netaddr := &net.UDPAddr{IP: net.ParseIP(strip), Port: int(p)}

		sock, err := net.ListenUDP("udp4", netaddr)
		if err != nil {
			fmt.Println("Listen Failure", err)
			return
		}


		usua.AddIP4Str(sock.LocalAddr().String())


		ul.sockarr = append(ul.sockarr,sock)
		ul.wg.Add(1)
		go ul.sockRecv(sock)
	}

	ul.listenAddr = usua

	fmt.Println("Server will start at:")
	usua.PrintAll()

	ul.wg.Wait()
	fmt.Println("Server Closed",)
}

func (ul *udplisten)GetListenAddr() address.UdpAddresser  {
	return ul.listenAddr
}

func (ul *udplisten)sockRecv(conn *net.UDPConn)  {
	defer ul.wg.Done()

	for{
		buf:=make([]byte,2048)
		nr,addr,err:=conn.ReadFromUDP(buf)
		if err!=nil{
			fmt.Println("sock Recv err",err.Error())
			return
		}
		//fmt.Println(addr.String())
		cp := NewConnPacket()
		if err=cp.DeSerialize(buf[0:nr]);err!=nil{
			fmt.Println("Deserialize error",err.Error())
			continue
		}
		cs:=GetConnStoreInstance()
		cs.Update(string(cp.GetUid()),conn,addr)


		if cp.GetTyp() == CONN_PACKET_TYP_KA{
			continue
		}

		c := cs.GetConn(string(cp.GetUid()))
		if c==nil{
			continue
		}
		c.Push(cp)
	}
}

func (ul *udplisten)Close()  {
	for _,sock:=range ul.sockarr{
		sock.Close()
	}
	ul.sockarr = nil
}