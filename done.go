package nbsnetwork

import (
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/translayer/recv"
	"github.com/kprc/nbsnetwork/file"
)

func NetWorkDone()  {

	store.GetStreamStoreInstance().Stop()
	store.GetBlockStoreInstance().Stop()
	file.GetFileStoreInstance().Stop()

	recv.ReceiveFromUdpStop()

	tools.GetNbsTickerInstance().Stop()


}