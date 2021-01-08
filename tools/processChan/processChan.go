package processChan

import (
	"io"
	"log"
	"net"
	"os"
	"path"
)

func ReceivePasswd(dir string) string {

	sockFile := path.Join(dir,"passwd.sock")

	if err:=os.RemoveAll(sockFile);err!=nil{
		log.Fatal("remove file failed")
		return ""
	}

	l,err:=net.Listen("unix",sockFile)
	if err!=nil{
		log.Fatal("listen error:", err)
	}

	defer l.Close()


	conn, err:=l.Accept()
	if err!=nil{
		log.Fatal("accept error: ",err)
	}
	defer conn.Close()
	if e,p:=readpaswd(conn);e!=nil{
		log.Fatal("conn error:",e.Error())
		return ""
	}else{
		return p
	}
}

func readpaswd(conn net.Conn) (err error,passwd string) {

	defer conn.Close()

	buf:=make([]byte,1024)
	n,err:=conn.Read(buf)
	if err!=nil && err!=io.EOF{
		return err,""
	}

	return nil,string(buf[:n])

}

func SendPasswd(dir string,passwd string)  {

	sockFile := path.Join(dir,"passwd.sock")

	c,err:=net.Dial("unix",sockFile)
	if err!=nil{
		log.Fatal("connect to uds failed")
	}
	defer c.Close()

	n,err:=c.Write([]byte(passwd))
	if err!=nil || n!=len(passwd){
		log.Fatal("write failed")
	}
}
