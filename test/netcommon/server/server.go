package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/tools"
	"strconv"
	"github.com/kprc/nbsnetwork/translayer/store"
	_ "github.com/kprc/nbsnetwork"
)

func main()  {

	server:=netcommon.GetUpdListenInstance()
	go server.Run("0.0.0.0",22113)
	//go server.Run("localaddress",0)
	cs:=netcommon.GetConnStoreInstance()

	for {
		rb := cs.Read()
		fmt.Println(string(rb.GetConnPacket().GetData()))
		s,p:=rb.GetUdpConn().GetAddr().FirstS()
		rb.GetUdpConn().Send([]byte("rcv ip: "+s+" port: "+strconv.Itoa(int(p))),store.UDP_MESSAGE)
	}

	tools.GetNbsTickerInstance().Stop()
}
