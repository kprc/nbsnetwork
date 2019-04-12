package file

import (
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
)

/**
uf UdpFile
f *os.File
 */

type filestoreblk struct {
	key store.UdpStreamKey
	blk interface{}
	lastAccessTime int64
}

func (fsb *filestoreblk)GetKey() store.UdpStreamKey {
	return fsb.key
}


type filestore struct {
	hashlist.HashList
	tick chan int64
}

type FileStore interface {
	AddFile(f interface{})
	DelFile(f interface{})
	FindFileDo(f interface{}, arg interface{}, do list.FDo)(r interface{}, err error)
	Run()
}



