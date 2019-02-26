package main

import (
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"github.com/kprc/nbsnetwork/netcommon"
	"time"
	"fmt"
)

func main()  {
	tick:=tools.GetNbsTickerInstance()
	var wg sync.WaitGroup
	wg.Add(1)
	go tick.Run(&wg)
	server:=netcommon.GetUpdListenInstance()
	go server.Run("0.0.0.0",22113)
	time.Sleep(time.Second*10)
	cs:=netcommon.GetConnStoreInstance()

	conn:=cs.First()

	rcv,_:=conn.Read()
	fmt.Println(string(rcv))
	conn.Send([]byte("good world"))

	wg.Wait()

}
