package file

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/pb/file"
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/translayer/stream"
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/common/constant"
	"os"
	"path/filepath"
	"github.com/kprc/nbsnetwork/rpc"
)

type udpfile struct {
	FileHead
	streamid uint64
	resume bool
	startSize int64
}

type UdpFile interface {
	FileHead
	SetStreamId(id uint64)
	GetStreamId() uint64
	SetResume(b bool)
	GetResume() bool
	SetStartSize(size int64)
	GetStartSize() int64
	Serialize() ([]byte,error)
	DeSerialize(data []byte)  error
	Clone() UdpFile
}

func NewUdpFile(fh FileHead) UdpFile  {
	return &udpfile{fh,0,false,0}
}

func NewEmptyUdpFile() UdpFile  {
	efh:=NewEmptyFileHead()
	return NewUdpFile(efh)
}

func (uf *udpfile)SetStreamId(id uint64)  {
	uf.streamid = id
}

func (uf *udpfile)GetStreamId() uint64  {
	return uf.streamid
}

func (uf *udpfile)SetResume(b bool)  {
	uf.resume = b
}

func (uf *udpfile)GetResume() bool  {
	return uf.resume
}

func (uf *udpfile)SetStartSize(size int64)  {
	uf.startSize = size
}

func (uf *udpfile)GetStartSize() int64  {
	return uf.startSize
}

func (uf *udpfile)Clone() UdpFile  {
    n:=NewEmptyUdpFile()
    n.SetSize(uf.GetSize())
    n.SetName(uf.GetName())
    n.SetPath(uf.GetPath())
    n.SetStrHash(uf.GetStrHash())
    n.SetStreamId(uf.GetStreamId())
    n.SetResume(uf.GetResume())
    n.SetStartSize(uf.GetStartSize())

	return n
}


func (uf *udpfile)Serialize() ([]byte,error)  {
	puf := &file.Udpfile{}

	puf.Name = []byte(uf.GetName())
	puf.Size = uf.GetSize()
	puf.Streamid = uf.GetStreamId()
	puf.Strhash = []byte(uf.GetStrHash())
	puf.Resume = uf.GetResume()
	puf.StartSize = uf.GetStartSize()

	return proto.Marshal(puf)
}

func (uf *udpfile)DeSerialize(data []byte)  error  {
	puf := &file.Udpfile{}

	if err:=proto.Unmarshal(data,puf);err!=nil{
		return err
	}
	uf.streamid = puf.GetStreamid()
	uf.SetName(string(puf.GetName()))
	uf.SetStrHash(string(puf.GetStrhash()))
	uf.SetSize(puf.GetSize())
	uf.SetResume(puf.GetResume())
	uf.SetStartSize(puf.GetStartSize())

	return nil
}

type filesend struct {
	UdpFile
	conn netcommon.UdpConn
	f *os.File
}

type UdpSend interface {
	UdpFile
	Send() error
	ResumeSend() error
	SetConn(conn netcommon.UdpConn)
}

func NewUdpFileSend(uf UdpFile,conn netcommon.UdpConn) UdpSend  {
	return &filesend{uf,conn,nil}
}

func (us *filesend)openFile() error {
	if fi,err:=os.Open(filepath.Join(us.GetPath(),us.GetName()));err!=nil{
		return err
	}else{
		us.f = fi
	}

	return nil

}

func (us *filesend)openFileAndSeek(size int64) error  {
	if err:=us.openFile();err!=nil{
		return err
	}

	if _,err:=us.f.Seek(size,os.SEEK_SET);err!=nil{
		return err
	}

	return nil
}


func (us *filesend)closeFile()  {
	if us.f !=nil{
		us.f.Close()
		us.f = nil
	}
}

func (us *filesend)Send() error  {
	ustream:=stream.NewUdpStream(us.conn,false)
	sid,err:=ustream.GetStreamId()
	if err!=nil{
		return err
	}
	us.SetStreamId(sid)

	msgsend:=message.NewReliableMsg(us.conn)
	msgsend.SetAppTyp(constant.FILE_DESC_HANDLE)
	var snddata []byte
	snddata,err=us.Serialize()
	if err!=nil {
		return err
	}
	err = msgsend.ReliableSend(snddata)
	if err!=nil{
		return err
	}

	ustream.SetAppTyp(constant.FILE_STREAM_HANDLE)
	ustream.SetTimeOut(30000)
	if us.f == nil{
		if err = us.openFile();err!=nil{
			return err
		}
	}
	err=ustream.ReliableSend(us.f)

	us.closeFile()

	return err
}

func (us *filesend)ResumeSend() error  {
	ustream:=stream.NewUdpStream(us.conn,false)

	if sid,err:=ustream.GetStreamId();err!=nil{
		return err
	}else{
		us.SetStreamId(sid)
	}

	us.SetResume(true)

	rpcsend := rpc.NewM2MRpc(us.conn,constant.FILE_DESC_HANDLE,nil,nil)

	var r interface{}
	var snddata []byte
	var err error

	if snddata,err=us.Serialize(); err!=nil{
		return err
	}
	if r,err = rpcsend.RPCSend(snddata);err!=nil{
		return err
	}

	uf:=r.(UdpFile)
	startSize:=uf.GetStartSize()

	ustream.SetAppTyp(constant.FILE_STREAM_HANDLE)
	ustream.SetTimeOut(30000)
	if us.f == nil{
		if err = us.openFileAndSeek(startSize);err!=nil{
			return err
		}
	}
	err=ustream.ReliableSend(us.f)

	us.closeFile()

	return err

}

func (us *filesend)SetConn(conn netcommon.UdpConn)  {
	us.conn = conn
}



