package main

import (
	_ "github.com/kprc/nbsnetwork"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork"
	"os"
	"github.com/kprc/nbsnetwork/file"
)

func main()  {

	if len(os.Args)>1{
		file.SetSaveFilePath(os.Args[1])
	}

	server:=netcommon.GetUpdListenInstance()
	server.Run("0.0.0.0",22113)

	server.Close()

	nbsnetwork.NetWorkDone()
}