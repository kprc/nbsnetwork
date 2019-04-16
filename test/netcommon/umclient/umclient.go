package main

import (
	_ "github.com/kprc/nbsnetwork"
	"os"
	"path/filepath"

	"fmt"

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

	return


	//uc:=netcommon.NewUdpCreateConnection(ips,"",22113,0)
	//
	//if err:=uc.Dial();err!=nil{
	//	fmt.Print("Dial Error",err.Error())
	//	return
	//}
	//uc.ConnSync()
	//go uc.Connect()
	//r:=uc.WaitConnReady()
	//
	//if !r{
	//	uc.Close()
	//	fmt.Println("Can't Connect to peer")
	//	return
	//}
	//
	//uf:=file.NewEmptyUdpFile()



}


