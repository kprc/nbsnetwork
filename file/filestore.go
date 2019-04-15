package file

import (
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"sync"
	"github.com/kprc/nbsnetwork/tools"

)

type filestoreblk struct {
	key store.UdpStreamKey
	blk interface{}
	lastAccessTime int64
	timeoutInterval int32
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
	AddFileWithParam(f interface{},timeoutInterval int32)
	DelFile(f interface{})
	FindFileDo(f interface{}, arg interface{}, do list.FDo)(r interface{}, err error)
	Run()
}

var (
	fsbInstance FileStore
	fsbInstanceLock sync.Mutex
	fsbAccessTime int64
	fsbTimeOutInterval int64 = 1000
)

func newFileStore() FileStore {
	hl:=hashlist.NewHashList(0x80,store.StreamKeyHash,store.StreamKeyEquals)
	fs:=&filestore{hl,make(chan int64,64)}

	t:=tools.GetNbsTickerInstance()
	t.Reg(&fs.tick)
	fsbAccessTime = tools.GetNowMsTime()

	return fs
}

func GetFileStoreInstance()  FileStore{
	if fsbInstance == nil{
		fsbInstanceLock.Lock()
		defer fsbInstanceLock.Unlock()
		if fsbInstance == nil{
			fsbInstance = newFileStore()
		}

	}

	return fsbInstance
}

func RefreshFSB(v interface{}){
	blk:=v.(*filestoreblk)

	blk.lastAccessTime = tools.GetNowMsTime()
}

func GetFileBlk()  {
	
}



func (fs *filestore)addFile(f interface{},timeoutInterval int32)  {
	fb:=f.(FileBlk)

	fsb:=&filestoreblk{blk:f}
	if timeoutInterval < 5000 {
		timeoutInterval = 5000
	}

	fsb.timeoutInterval = timeoutInterval

	fsb.key = fb.GetKey()

	fsb.lastAccessTime = tools.GetNowMsTime()

	fs.Add(fsb)
}

func (fs *filestore)AddFile(f interface{})  {
	fs.addFile(f,0)
}

func (fs *filestore)AddFileWithParam(f interface{},timeoutInterval int32)  {
	fs.addFile(f,timeoutInterval)
}

func (fs *filestore)DelFile(f interface{})  {
	fs.Del(f)
}

func (fs *filestore)FindFileDo(f interface{}, arg interface{}, do list.FDo)(r interface{}, err error)  {
	return fs.FindDo(f,arg,do)
}


func (fs *filestore)doTimeOut()  {
	type fs2del struct {
		arrdel []*filestoreblk
	}

	arr2del := &fs2del{arrdel:make([]*filestoreblk,0)}

	fdo:= func(arg interface{}, v interface{})(r interface{},err error) {
		fsb:=v.(*filestoreblk)
		l:=arg.(*fs2del)

		curtime:=tools.GetNowMsTime()
		tv:=curtime - fsb.lastAccessTime
		if tv > int64(fsb.timeoutInterval){
			l.arrdel = append(l.arrdel,fsb)
			CloseFile(fsb.blk)
		}

		return
	}

	fs.TraversAll(arr2del,fdo)

	for _,fsb:=range arr2del.arrdel{
		fs.DelFile(fsb)
	}
}

func (fs *filestore)Run()  {
	select {
	case <-fs.tick:
		if tools.GetNowMsTime() - fsbAccessTime > fsbTimeOutInterval{
			fsbAccessTime = tools.GetNowMsTime()
			fs.doTimeOut()
		}
	}
}


