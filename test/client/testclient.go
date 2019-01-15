package main

import (
	"github.com/kprc/nbsnetwork/client"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/send"
)

func main()  {

	bs := send.GetInstance()
	bs.TimeOut()

	c := client.NewUdpClient("192.168.107.242","",11223,0)
	c.Dial()

	c.SendBytes([]byte("Title Ping")[:],constant.MSG_PING,[]byte("PING TO SERVER")[:])

}


