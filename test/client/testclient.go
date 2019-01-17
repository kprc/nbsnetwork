package main

import (
	"github.com/kprc/nbsnetwork/client"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/recv"
	"github.com/kprc/nbsnetwork/send"
	"github.com/kprc/nbsnetwork/test/msghandle"
)

func main()  {

	msghandle.RegPingAckMsg()

	bs := send.GetInstance()
	go bs.TimeOut()
	rmr:=recv.GetInstance()
	go rmr.TimeOut()

	//ip:="192.168.107.242"
	ip:="192.168.103.66"

	c := client.NewUdpClient(ip,"",11223,0)
	c.Dial()

	c.SendBytes([]byte("Title Ping")[:],constant.MSG_PING,[]byte("PING TO SERVER")[:])

}


