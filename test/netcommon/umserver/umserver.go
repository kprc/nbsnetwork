package main

import (
	_ "github.com/kprc/nbsnetwork"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/netcommon"
)

func main()  {

	server:=netcommon.GetUpdListenInstance()
	server.Run("0.0.0.0",22113)

	tools.GetNbsTickerInstance().Stop()
}