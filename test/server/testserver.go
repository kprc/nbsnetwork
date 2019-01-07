package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"github.com/kprc/nbsnetwork/rw"
	"github.com/kprc/nbsnetwork/server"
	"io"
)

func main()  {
	fmt.Println("Test Server")

	ws:= getPingWS(nil)

	ws.Write([]byte("hello,"))
	ws.Write([]byte("world"))


	ws.(rw.UdpBytesWriterSeeker).PrintAll()

	RegMsg()


	us := server.GetUdpServer()

	us.Run("",11223)


	fmt.Println("End Server")
}

func RegMsg()  {
	mh:=regcenter.NewMsgHandler()

	mh.SetHandler(handlePing)
	mh.SetWSNew(getPingWS)

	mi:=regcenter.GetMsgCenterInstance()

	mi.AddHandler(constant.MSG_PING,mh)

}

func getPingWS(param interface{}) io.WriteSeeker {
	ws:= rw.NewWriteSeeker(param)

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