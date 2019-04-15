package nbsnetwork

import (
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/file"
)

func init()  {
	tick:=tools.GetNbsTickerInstance()
	go tick.Run()

	file.FileRegister()
}
