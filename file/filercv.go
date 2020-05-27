package file

import (
	"fmt"
	"github.com/kprc/nbsnetwork/applayer"
	"github.com/kprc/nbsnetwork/bus"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/rpc"
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
)

func FileRegister() {
	fmt.Println("File Handle Register")

	abs := applayer.GetAppBlockStore()

	abs.Reg(constant.FILE_DESC_HANDLE, handleFileHead)
	abs.Reg(constant.FILE_STREAM_HANDLE, handleFileStream)
	abs.Reg(constant.FILE_START_SIZE, handleFileStartSize)

}

func sendResumeDesc(v interface{}, uid string) {
	conn := netcommon.GetConnStoreInstance().GetConn(uid)
	if conn == nil || !conn.Status() {
		return
	}
	uf := v.(UdpFile)
	fo := NewFileOp(nil)
	size := fo.GetFileSize(uf.GetFileName())

	uf.SetStartSize(size)
	mr := message.NewReliableMsg(conn)
	mr.SetAppTyp(constant.FILE_START_SIZE)
	if snddata, err := uf.Serialize(); err != nil {
		return
	} else {
		err1 := mr.ReliableSend(snddata)
		if err1 != nil {
			fmt.Println(err1.Error())
		}
	}
}

func handleFileStartSize(rcv interface{}, arg interface{}) (v interface{}, err error) {
	cb := rcv.(applayer.CtrlBlk)
	um := cb.GetUdpMsg()
	uf := NewEmptyUdpFile()
	err = uf.DeSerialize(um.GetData())
	if err != nil {
		return nil, err
	}
	sn := uf.GetStreamId()
	bk := store.NewBlockKey(sn)
	fdo := func(arg interface{}, v interface{}) (r interface{}, err error) {
		rb := rpc.GetRpcBlock(v)

		rb.GetRpcDo()(rb.GetData(), arg, rb.GetResponseChan(), false)

		return
	}

	rpcstore := rpc.GetRpcStore()
	return rpcstore.FindRpcBlockDo(bk, uf, fdo)
}

func handleFileHead(rcv interface{}, arg interface{}) (v interface{}, err error) {
	cb := rcv.(applayer.CtrlBlk)
	um := cb.GetUdpMsg()

	uf := NewEmptyUdpFile()

	err = uf.DeSerialize(um.GetData())
	if err != nil {
		return nil, err
	}

	uid := cb.GetRcvBlk().GetConnPacket().GetUid()
	streamid := uf.GetStreamId()

	key := store.NewUdpStreamKeyWithParam(string(uid), streamid)

	fs := GetFileStoreInstance()

	if !findFileBlk(key) {
		fb := NewFileBlk()
		fb.SetKey(key)
		fb.SetUdpFile(uf)
		fs.AddFileWithParam(fb, int32(constant.FILE_STORE_TIMEOUT))
	}

	if uf.GetResume() {
		bb := bus.NewBusBlock()
		bb.SetUid(string(uid))
		n := uf.Clone()
		n.SetStreamId(um.GetSn())
		bb.SetData(n)
		bb.SetFSendData(sendResumeDesc)
		bs := bus.GetBusStoreInstance()
		bs.AddBlock(bb)
	}

	return nil, nil
}

func findFileBlk(key store.UdpStreamKey) bool {
	fdo := func(arg interface{}, v interface{}) (ret interface{}, err error) {
		RefreshFSB(v)
		return v, nil
	}

	fs := GetFileStoreInstance()

	if _, err := fs.FindFileDo(key, nil, fdo); err != nil {
		return false
	}

	return true
}

func openFile(key store.UdpStreamKey) (io.WriteCloser, error) {
	fdo := func(arg interface{}, v interface{}) (ret interface{}, err error) {
		blk := GetFileBlk(v).(FileBlk)
		filename := blk.GetUdpFile().GetFileName()
		if blk.GetFileOp() == nil {
			fo := NewFileOp(nil)
			blk.SetFileOp(fo)
		}
		mode := NEW_CREATE
		if blk.GetUdpFile().GetResume() {
			mode = APPEND_CREATE
		}
		blk.GetFileOp().CreateFile(filename, mode)

		return blk.GetFileOp(), nil
	}

	fs := GetFileStoreInstance()
	if wc, err := fs.FindFileDo(key, nil, fdo); err != nil {
		return nil, err
	} else {
		return wc.(FileOp), nil
	}
}

func freshFile(key store.UdpStreamKey) error {
	fdo := func(arg interface{}, v interface{}) (ret interface{}, err error) {
		RefreshFSB(v)
		return
	}

	fs := GetFileStoreInstance()

	fs.FindFileDo(key, nil, fdo)

	return nil
}

func closeFile(key store.UdpStreamKey) error {
	fmt.Println("close File from fdo")
	fdo := func(arg interface{}, v interface{}) (ret interface{}, err error) {
		blk := GetFileBlk(v).(FileBlk)
		if blk.GetFileOp() != nil {
			blk.GetFileOp().Close()
		}

		return v, nil
	}

	fs := GetFileStoreInstance()
	if fsb, err := fs.FindFileDo(key, nil, fdo); err != nil {
		return err
	} else {
		if fsb != nil {
			fs.DelFile(fsb)
		}
	}

	return nil
}

func handleFileStream(rcv interface{}, arg interface{}) (v interface{}, err error) {
	cb := rcv.(applayer.CtrlBlk)
	uid := cb.GetRcvBlk().GetConnPacket().GetUid()
	sn := cb.GetUdpMsg().GetSn()

	key := store.NewUdpStreamKeyWithParam(string(uid), sn)

	h := arg.(int)

	if h == constant.OPEN_FILE {
		return openFile(key)
	} else if h == constant.CLOSE_FILE {
		closeFile(key)
	} else if h == constant.REFRESH_FILE {
		freshFile(key)
	}

	return
}
