package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/tools"
)

func main()  {
	tick:=tools.GetNbsTickerInstance()
	go tick.Run()
	server:=netcommon.GetUpdListenInstance()
	go server.Run("0.0.0.0",22113)
	//go server.Run("localaddress",0)
	cs:=netcommon.GetConnStoreInstance()

	for {
		rb := cs.Read()
		fmt.Println(string(rb.GetConnPacket().GetData()))
		rb.GetUdpConn().Send([]byte("good world"))
	}

	tick.Stop()
}
