package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

func main()  {
	tick:=tools.GetNbsTickerInstance()
	var wg sync.WaitGroup
	wg.Add(1)
	go tick.Run(&wg)
	server:=netcommon.GetUpdListenInstance()
	go server.Run("0.0.0.0",22113)

	cs:=netcommon.GetConnStoreInstance()

	for {
		rb := cs.Read()
		fmt.Println(string(rb.GetConnPacket().GetData()))
		rb.GetUdpConn().Send([]byte("good world"))
	}

	wg.Wait()

}
