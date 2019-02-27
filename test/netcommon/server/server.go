package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"time"
)

func main()  {
	tick:=tools.GetNbsTickerInstance()
	var wg sync.WaitGroup
	wg.Add(1)
	go tick.Run(&wg)
	server:=netcommon.GetUpdListenInstance()
	go server.Run("0.0.0.0",22113)

	cs:=netcommon.GetConnStoreInstance()

	var conn netcommon.UdpConn

	for  {
		conn=cs.First()
		time.Sleep(time.Second*1)
		if conn !=nil {
			break
		}
	}



	rcv,_:=conn.Read()
	fmt.Println(string(rcv))
	conn.Send([]byte("good world"))

	wg.Wait()

}
