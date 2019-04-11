package file

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"io"
	"github.com/kprc/nbsnetwork/pb/file"
	"github.com/gogo/protobuf/proto"
)

type udpfile struct {
	FileHead
	streamid uint64
}

type UdpFile interface {
	FileHead
	SetStreamId(id uint64)
	GetStreamId() uint64
	Serialize() ([]byte,error)
	DeSerialize(data []byte)  error
}

func NewUdpFile(fh FileHead) UdpFile  {
	return &udpfile{fh,0}
}

func (uf *udpfile)SetStreamId(id uint64)  {
	uf.streamid = id
}

func (uf *udpfile)GetStreamId() uint64  {
	return uf.streamid
}

func (uf *udpfile)Serialize() ([]byte,error)  {
	puf := &file.Udpfile{}

	puf.Name = []byte(uf.GetName())
	puf.Size = uf.GetSize()
	puf.Streamid = uf.GetStreamId()
	puf.Strhash = []byte(uf.GetStrHash())

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

	return nil
}


type filesend struct {
	conn netcommon.UdpConn
	io.Reader
	UdpFile
}


type UdpSend interface {
	UdpFile
	Send() error
	SetConn(conn netcommon.UdpConn)
	SetFileHead(fh UdpFile)
}

func (us *filesend)Send() error  {

	
	
	return nil
}

func (us *filesend)SetConn(conn netcommon.UdpConn)  {
	us.conn = conn
}

func (us *filesend)SetFileHead(fh UdpFile)  {
	us.UdpFile = fh
}



