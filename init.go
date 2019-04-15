package nbsnetwork

import (
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/file"
	"github.com/kprc/nbsnetwork/translayer/store"
)

func init()  {
	tick:=tools.GetNbsTickerInstance()
	go tick.Run()

	file.FileRegister()


	runstore()
}


func runstore(){
	fs:=file.GetFileStoreInstance()
	go fs.Run()

	msgstore:=store.GetBlockStoreInstance()
	go msgstore.Run()

	ss:=store.GetStreamStoreInstance()

	go ss.Run()

}