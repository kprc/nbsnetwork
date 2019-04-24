package file

import (
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
	"fmt"
	"github.com/kprc/nbsnetwork/applayer"
)



func FileRegister()  {
	fmt.Println("File Handle Register")

	abs:=applayer.GetAppBlockStore()

	abs.Reg(constant.FILE_DESC_HANDLE,handleFileHead)
	abs.Reg(constant.FILE_STREAM_HANDLE,handleFileStream)

}

func handleFileHead(rcv interface{},arg interface{}) (v interface{},err error)  {
	cb:=rcv.(applayer.CtrlBlk)
	um:=cb.GetUdpMsg()

	uf:=NewEmptyUdpFile()

	err=uf.DeSerialize(um.GetData())
	if err!=nil{
		return nil,err
	}

	uid:=cb.GetRcvBlk().GetConnPacket().GetUid()
	streamid:=uf.GetStreamId()

	key:=store.NewUdpStreamKeyWithParam(string(uid),streamid)

	fb:=NewFileBlk()
	fb.SetKey(key)
	fb.SetUdpFile(uf)

	fs:=GetFileStoreInstance()

	if !findFileBlk(key) {
		fs.AddFileWithParam(fb,int32(constant.FILE_STORE_TIMEOUT))
	}

	return nil,nil
}

func findFileBlk(key store.UdpStreamKey) bool  {
	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error){
		RefreshFSB(v)
		return v,nil
	}

	fs:=GetFileStoreInstance()

	if _,err:=fs.FindFileDo(key,nil,fdo);err!=nil{
		return false
	}

	return true
}

func openFile(key store.UdpStreamKey) (io.WriteCloser,error) {
	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error) {
		blk:=GetFileBlk(v).(FileBlk)
		filename := blk.GetUdpFile().GetFileName()
		if blk.GetFileOp() == nil{
			fo:=NewFileOp(nil)
			blk.SetFileOp(fo)
		}
		blk.GetFileOp().CreateFile(filename,NEW_CREATE)

		return blk.GetFileOp(),nil
	}

	fs:=GetFileStoreInstance()
	if wc,err:=fs.FindFileDo(key,nil,fdo);err!=nil{
		return nil,err
	}else{
		return wc.(FileOp),nil
	}
}


func openFileContinue(key store.UdpStreamKey) (io.WriteCloser,error) {
	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error) {
		blk:=GetFileBlk(v).(FileBlk)
		filename := blk.GetUdpFile().GetFileName()
		if blk.GetFileOp() == nil{
			fo:=NewFileOp(nil)
			blk.SetFileOp(fo)
		}
		blk.GetFileOp().CreateFile(filename,APPEND_CREATE)

		return blk.GetFileOp(),nil
	}

	fs:=GetFileStoreInstance()
	if wc,err:=fs.FindFileDo(key,nil,fdo);err!=nil{
		return nil,err
	}else{
		return wc.(FileOp),nil
	}
}

func freshFile(key store.UdpStreamKey) error  {
	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error) {
		RefreshFSB(v)
		return
	}

	fs:=GetFileStoreInstance()

	fs.FindFileDo(key,nil,fdo)

	return nil
}

func closeFile(key store.UdpStreamKey) error {
	fmt.Println("close File from fdo")
	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error) {
		blk:=GetFileBlk(v).(FileBlk)
		if blk.GetFileOp() != nil {
			blk.GetFileOp().Close()
		}

		return v,nil
	}

	fs:=GetFileStoreInstance()
	if fsb,err:=fs.FindFileDo(key,nil,fdo);err!=nil{
		return err
	}else {
		if fsb !=nil{
			fs.DelFile(fsb)
		}
	}

	return nil
}




func handleFileStream(rcv interface{},arg interface{}) (v interface{},err error) {
	cb:=rcv.(applayer.CtrlBlk)
	uid:=cb.GetRcvBlk().GetConnPacket().GetUid()
	sn:=cb.GetUdpMsg().GetSn()

	key:=store.NewUdpStreamKeyWithParam(string(uid),sn)

	h:=arg.(int)

	if h == constant.OPEN_FILE {
		return openFile(key)
	}else if h== constant.CLOSE_FILE{
		closeFile(key)
	}else if h == constant.REFRESH_FILE{
		freshFile(key)
	}else if h == constant.OPEN_FILE_CONTINUE{
		openFileContinue(key)
	}

	return
}








