package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/test/msghandle"
	"io"
	"github.com/kprc/nbsnetwork/server"
)

func main()  {
	fmt.Println("Test Server")

	//testWriteSeeker()
	//testReadSeeker()

	msghandle.RegPingMsg()

	us := server.GetUdpServer()

	us.Run("",11223)

	fmt.Println("End Server")
}



func testWriteSeeker()  {
	ws:= regcenter.GetMsgCenterInstance().GetHandler(constant.MSG_PING).GetWSNew()(nil)

	ws.Write([]byte("hello,"))
	ws.Write([]byte("world"))


	ws.(netcommon.UdpBytesWriterSeeker).PrintAll()

	fmt.Println(string(ws.(netcommon.UdpBytesWriterSeeker).GetBytes()))
}

func testReadSeeker()  {
	rs:=netcommon.NewReadSeeker([]byte("hello,world"))

	rs.PrintAll()

	p:=make([]byte,5)

	for  {
		nr,err := rs.Read(p)
		if nr>0 {
			fmt.Println(string(p[:nr]))
		}
		if err == io.EOF || err != nil{
			break
		}
	}

}
