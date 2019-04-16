package main

import (
	_ "github.com/kprc/nbsnetwork"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork"
)

func main()  {

	server:=netcommon.GetUpdListenInstance()
	server.Run("0.0.0.0",22113)

	nbsnetwork.NetWorkDone()
}