package main

import (
	"github.com/kprc/nbsnetwork"
	_ "github.com/kprc/nbsnetwork"
	"github.com/kprc/nbsnetwork/file"
	"github.com/kprc/nbsnetwork/netcommon"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		file.SetSaveFilePath(os.Args[1])
	}

	server := netcommon.GetUpdListenInstance()
	server.Run("0.0.0.0", 22113)

	server.Close()

	nbsnetwork.NetWorkDone()
}
