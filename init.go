package nbsnetwork

import (
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/file"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/translayer/recv"
)

func init()  {
	tick:=tools.GetNbsTickerInstance()
	go tick.Run()

	file.FileRegister()

	runstore()

	go recv.ReceiveFromUdpConn()
}


func runstore(){
	fs:=file.GetFileStoreInstance()
	go fs.Run()

	msgstore:=store.GetBlockStoreInstance()
	go msgstore.Run()

	ss:=store.GetStreamStoreInstance()

	go ss.Run()

}