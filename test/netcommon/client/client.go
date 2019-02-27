package main

import (
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"github.com/kprc/nbsnetwork/netcommon"
	"fmt"
	"time"
)

func main()  {
	tick:=tools.GetNbsTickerInstance()
	var wg sync.WaitGroup
	wg.Add(1)
	go tick.Run(&wg)
	uc:=netcommon.NewUdpCreateConnection("127.0.0.1","",22113,0)
	uc.Dial()
	go uc.Connect()

	time.Sleep(time.Second*2)
	uc.Send([]byte("hello world"))

	cs:=netcommon.GetConnStoreInstance()

	icnt:=0

	for {
		c := cs.Read()
		fmt.Println(string(c.GetConnPacket().GetData()))
		uc.Send([]byte("client send time:"+time.Now().String()))
		icnt++
		if icnt>10{
			break
		}
	}
	uc.Close()
	tick.Stop()
	wg.Wait()

}
