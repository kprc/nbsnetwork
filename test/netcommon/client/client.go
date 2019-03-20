package main

import (
	"fmt"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/netcommon"
	tm "github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/tools"
	"time"
)

func main()  {
	tick:=tools.GetNbsTickerInstance()
	go tick.Run()
	//uc:=netcommon.NewUdpCreateConnection("127.0.0.1","",22113,0)
	uc:=netcommon.NewUdpCreateConnection("127.0.0.1","",0,0)
	uc.Dial()
	uc.Hello()
	go uc.Connect()
	r:=uc.WaitHello()

	if !r{
		uc.Close()
		fmt.Println("Can't Connect to peer")
		return
	}

	uc.Send([]byte("hello world"))

	cs:=netcommon.GetConnStoreInstance()
	cs.Add(nbsid.GetLocalId().String(),uc)

	icnt:=0

	for {
		c, err := cs.ReadAsync()
		if c != nil && err==nil {
			fmt.Println(string(c.GetConnPacket().GetData()))
		}
		err = uc.Send([]byte("client send time:" + time.Now().String()))
		if err != nil {
			fmt.Println(err.(nbserr.NbsErr).Errmsg)
			break
		}
		time.Sleep(time.Second*1)
		icnt++
		if icnt>15{
			break
		}
	}
	uc.Close()
	tick.Stop()
}

func read(cs netcommon.ConnStore,quit *chan int)  {
	for {
		c := cs.Read()

		um:=tm.NewUdpMsg(0,nil)
		if err:=um.DeSerialize(c.GetConnPacket().GetData());err!=nil{
			continue
		}
		fmt.Println("um data:",string(um.GetData()))

		select {
		    case <-*quit:
				return
		    default:
		}
	}

}
