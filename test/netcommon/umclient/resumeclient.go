package main

import (
	_ "github.com/kprc/nbsnetwork"
	"os"
	"path/filepath"
	"fmt"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/file"
	"github.com/kprc/nbsnetwork"
	"github.com/pkg/errors"
	"github.com/kprc/nbsdht/dht/nbsid"
)

func main()  {
	ips:="127.0.0.1"
	//port:=22113
	userpath:=""
	if len(os.Args) > 1 {
		ips = os.Args[1]
		//if len(os.Args) > 2 {
		//	port,_ = strconv.Atoi(os.Args[2])
		//}
		if len(os.Args) > 2{
			userpath = os.Args[2]
		}
	}

	abspath,_:= filepath.Abs(userpath)
	fmt.Println(abspath,ips)

	fmt.Println(filepath.Dir(abspath),filepath.Base(abspath))

	nid:=nbsid.GetLocalId()
	fmt.Println("localid:",nid.String())

	for{
		if err:=sendfile(ips,abspath);err!=nil{
			fmt.Println(err.Error(),"Begin to Resume Send")
			continue
		}else {
			break
		}
	}

	nbsnetwork.NetWorkDone()
}

func sendfile(ips string, abspath string) error  {
	uc:=netcommon.NewUdpCreateConnection(ips,"",22113,0)

	if err:=uc.Dial();err!=nil{
		fmt.Print("Dial Error",err.Error())
		return err
	}
	uc.ConnSync()
	go uc.Connect()
	r:=uc.WaitConnReady()

	if !r{
		uc.Close()
		fmt.Println("Can't Connect to peer")
		return errors.New("Connect to peer error")
	}

	uf:=file.NewEmptyUdpFile()
	uf.SetSize(1024)
	uf.SetStrHash("testhash")
	uf.SetName(filepath.Base(abspath))
	uf.SetPath(filepath.Dir(abspath))

	ufs := file.NewUdpFileSend(uf,uc)

	err:=ufs.ResumeSend()
	if err!=nil{
		fmt.Println(err.Error())
		return err
	}

	return nil
}