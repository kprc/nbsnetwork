package nbsnetwork

import "github.com/kprc/nbsnetwork/tools"

func init()  {
	tick:=tools.GetNbsTickerInstance()
	go tick.Run()
}
