package main

import (
	"fmt"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/tools"
	"time"
)

func main()  {
	tick:=tools.GetNbsTickerInstance()
	go tick.Run()
	uc:=netcommon.NewUdpCreateConnection("192.168.107.242","",64099,0)
	uc.Dial()
	uc.Hello()
	go uc.Connect()
	uc.WaitHello()

	uc.Send([]byte("hello world"))

	cs:=netcommon.GetConnStoreInstance()
	cs.Add(nbsid.GetLocalId().String(),uc)

	icnt:=0

	for {
		c := cs.Read()
		fmt.Println(string(c.GetConnPacket().GetData()))
		c.GetUdpConn().Send([]byte("client send time:"+time.Now().String()))
		icnt++
		if icnt>10{
			break
		}
	}
	uc.Close()
	tick.Stop()
}
