package nbsnetwork

import (
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/translayer/recv"
	"github.com/kprc/nbsnetwork/file"
	"github.com/kprc/nbsnetwork/bus"
	"github.com/kprc/nbsnetwork/rpc"
)

func NetWorkDone()  {

	store.GetStreamStoreInstance().Stop()
	store.GetBlockStoreInstance().Stop()
	file.GetFileStoreInstance().Stop()
	bus.GetBusStoreInstance().Stop()
	rpc.GetRpcStore().Stop()

	recv.ReceiveFromUdpStop()

	tools.GetNbsTickerInstance().Stop()
}