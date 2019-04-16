package file

import (
	"github.com/kprc/nbsnetwork/applayer"
	"github.com/kprc/nbsnetwork/common/constant"

	"github.com/kprc/nbsnetwork/translayer/store"
)

func FileRegister()  {
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
		fs.AddFile(fb)
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

func openFile(key store.UdpStreamKey) error {
	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error) {
		blk:=GetFileBlk(v).(FileBlk)
		filename := blk.GetUdpFile().GetFileName()
		blk.GetFileOp().OpenFile(filename)

		return nil,nil
	}

	fs:=GetFileStoreInstance()
	if _,err:=fs.FindFileDo(key,nil,fdo);err!=nil{
		return err
	}

	return nil
}

func closeFile(key store.UdpStreamKey) error {
	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error) {
		blk:=GetFileBlk(v).(FileBlk)
		//filename := blk.GetUdpFile().GetFileName()
		//blk.GetFileOp().OpenFile(filename)
		blk.GetFileOp().Close()

		return nil,nil
	}

	fs:=GetFileStoreInstance()
	if _,err:=fs.FindFileDo(key,nil,fdo);err!=nil{
		return err
	}

	return nil
}




func handleFileStream(rcv interface{},arg interface{}) (v interface{},err error) {
	cb:=rcv.(applayer.CtrlBlk)
	uid:=cb.GetRcvBlk().GetConnPacket().GetUid()
	sn:=cb.GetUdpMsg().GetSn()

	key:=store.NewUdpStreamKeyWithParam(string(uid),sn)

	closeflag := arg.(bool)

	if !closeflag{
		openFile(key)
	}else {
		closeFile(key)
	}

	return
}








