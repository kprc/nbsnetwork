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
	uc:=netcommon.NewUdpCreateConnection("192.168.103.66","",22113,0)
	uc.Dial()
	go uc.Connect()

	time.Sleep(time.Second*2)
	uc.Send([]byte("hello world"))

	rcv,err:=uc.Read()

	if err!=nil{
		return
	}

	fmt.Println(string(rcv))
	uc.Close()
	wg.Wait()

}
