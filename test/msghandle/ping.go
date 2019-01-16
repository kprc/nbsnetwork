package msghandle

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"github.com/kprc/nbsnetwork/netcommon"
	"io"
)



func RegPingMsg()  {
	mh:=regcenter.NewMsgHandler()

	mh.SetHandler(handlePing)
	mh.SetWSNew(getPingWS)

	mi:=regcenter.GetMsgCenterInstance()

	mi.AddHandler(constant.MSG_PING,mh)

}

func getPingWS(param interface{}) io.WriteSeeker {
	ws:= netcommon.NewWriteSeeker(param)

	return ws

}

func handlePing(head interface{},data interface{},snd io.Writer) error  {

	if head != nil{
		sh:=head.([]byte)
		fmt.Println("Head is :",string(sh))
	}

	if data !=nil {
		sd := data.([]byte)
		fmt.Println("Data is :",string(sd))
	}

	snd.Write([]byte("Pong"))

	fmt.Println("Send Pong")

	return nil

}