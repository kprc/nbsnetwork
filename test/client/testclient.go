package main

import (
	"github.com/kprc/nbsnetwork/client"
	"github.com/kprc/nbsnetwork/common/constant"
)

func main()  {
	c := client.NewUdpClient("192.168.107.242","",11223,0)
	c.Dial()

	c.SendBytes([]byte("Title Ping")[:],constant.MSG_PING,"10366servber",[]byte("PING TO SERVER")[:])


}


