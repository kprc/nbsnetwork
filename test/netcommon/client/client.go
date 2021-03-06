package main

import (
	"fmt"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	_ "github.com/kprc/nbsnetwork"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/translayer/store"
	"os"
	"strconv"
	"time"
)

func main() {
	ips := "127.0.0.1"
	port := 22113
	if len(os.Args) > 1 {
		ips = os.Args[1]
		if len(os.Args) > 2 {
			port, _ = strconv.Atoi(os.Args[2])
		}
	}
	fmt.Println("Connectting to IP:", ips, " Port:", port)

	uc := netcommon.NewUdpCreateConnection(ips, "", uint16(port), 0)

	if err := uc.Dial(); err != nil {
		fmt.Print("Dial Error", err.Error())
		return
	}
	uc.ConnSync()
	go uc.Connect()
	r := uc.WaitConnReady()

	if !r {
		uc.Close()
		fmt.Println("Can't Connect to peer")
		return
	}

	uc.Send([]byte("hello world"), store.UDP_MESSAGE)

	cs := netcommon.GetConnStoreInstance()
	cs.Add(nbsid.GetLocalId().String(), uc)

	icnt := 0

	for {
		c, err := cs.ReadAsync()
		if c != nil && err == nil {
			fmt.Println(string(c.GetConnPacket().GetData()))
		}
		err = uc.Send([]byte("client send time:"+time.Now().String()), store.UDP_MESSAGE)
		if err != nil {
			fmt.Println(err.(nbserr.NbsErr).Errmsg)
			break
		}
		time.Sleep(time.Second * 1)
		icnt++
		if icnt > 15 {
			break
		}
	}
	uc.Close()
	tools.GetNbsTickerInstance().Stop()
	//tick.Stop()
}

//func read(cs netcommon.ConnStore,quit *chan int)  {
//	for {
//		c := cs.Read()
//
//		um:=tm.NewUdpMsg(0,nil)
//		if err:=um.DeSerialize(c.GetConnPacket().GetData());err!=nil{
//			continue
//		}
//		fmt.Println("um data:",string(um.GetData()))
//
//		select {
//		    case <-*quit:
//				return
//		    default:
//		}
//	}
//
//}
