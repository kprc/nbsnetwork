package file

import (
	"github.com/gogo/protobuf/io"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/netcommon"
	"sync"
)


type udptransfile struct {
	fd filedesc
	localpath string
	r io.Reader

	llock sync.Mutex
	lblock list.List
}


type UdpTransFile interface {
	SendFileDesc() error
	SendFile() error
}

func (utf *udptransfile)SendFileDesc(conn netcommon.UdpConn) error{
	
}

func (utf *udptransfile)SendFile() error{

}
