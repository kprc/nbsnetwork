package file

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"io"
)

type udpfile struct {
	conn netcommon.UdpConn
	fh FileHead
	io.Reader
	session uint64
}


type UdpFile interface {
	Send() error
	SetConn(conn netcommon.UdpConn)
	SetFileHead(fh FileHead)
}

func (uf *udpfile)Send() error  {
	//todo...
	
	
	return nil
}

func (uf *udpfile)SetConn(conn netcommon.UdpConn)  {
	uf.conn = conn
}

func (uf *udpfile)SetFileHead(fh FileHead)  {
	uf.fh = fh
}








