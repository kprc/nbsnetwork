package file

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
	"os"
	"sync"
)

const (
	APPEND_CREATE int = 1
	NEW_CREATE    int = 2
)

type fileop struct {
	f    *os.File
	lock sync.Mutex
}

type FileOp interface {
	io.Reader
	io.Writer
	io.Closer
	SetFile(f *os.File)
	IsClosed() bool
	CreateFile(filename string, mode int) error
	GetFileSize(filename string) int64
}

var (
	filenotseterr = nbserr.NbsErr{ErrId: nbserr.FILE_NOT_SET, Errmsg: "File Not Set"}
)

func NewFileOp(f *os.File) FileOp {
	return &fileop{f: f}
}

func (fo *fileop) GetFileSize(filename string) int64 {
	fi, err := os.Stat(filename)

	if err != nil {
		return 0
	}

	return fi.Size()
}

func (fo *fileop) CreateFile(filename string, mode int) error {
	if fo.f == nil {
		fo.lock.Lock()
		defer fo.lock.Unlock()
		if fo.f == nil {
			m := os.O_CREATE | os.O_WRONLY
			if mode == APPEND_CREATE {
				m |= os.O_APPEND
			}
			if f, err := os.OpenFile(filename, m, 0666); err != nil {
				return err
			} else {
				fo.f = f
			}
		}
	}
	return nil
}

func (fo *fileop) SetFile(f *os.File) {
	fo.lock.Lock()
	defer fo.lock.Unlock()
	fo.f = f
}

func (fo *fileop) Read(buf []byte) (int, error) {
	fo.lock.Lock()
	defer fo.lock.Unlock()
	if fo.f == nil {
		return 0, filenotseterr
	}

	return fo.f.Read(buf)
}

func (fo *fileop) Write(p []byte) (n int, err error) {
	fo.lock.Lock()
	defer fo.lock.Unlock()
	if fo.f == nil {
		return 0, filenotseterr
	}
	return fo.f.Write(p)
}

func (fo *fileop) Close() error {
	fo.lock.Lock()
	defer fo.lock.Unlock()
	if fo.f == nil {
		return filenotseterr
	}

	err := fo.f.Close()
	fo.f = nil

	return err
}

func (fo *fileop) IsClosed() bool {
	fo.lock.Lock()
	defer fo.lock.Unlock()
	if fo.f == nil {
		return true
	} else {
		return false
	}
}

type fileblk struct {
	uf  UdpFile
	fo  FileOp
	key store.UdpStreamKey
}

func CloseFile(v interface{}) {
	fb := v.(FileBlk)
	if fb.GetFileOp() != nil {
		fb.GetFileOp().Close()
	}
}

type FileBlk interface {
	SetUdpFile(uf UdpFile)
	GetUdpFile() UdpFile
	SetFileOp(op FileOp)
	GetFileOp() FileOp
	SetKey(key store.UdpStreamKey)
	GetKey() store.UdpStreamKey
}

func NewFileBlk() FileBlk {
	return &fileblk{}
}

func (fb *fileblk) SetUdpFile(uf UdpFile) {
	fb.uf = uf
}

func (fb *fileblk) GetUdpFile() UdpFile {
	return fb.uf
}

func (fb *fileblk) SetFileOp(op FileOp) {
	fb.fo = op
}

func (fb *fileblk) GetFileOp() FileOp {
	return fb.fo
}

func (fb *fileblk) SetKey(key store.UdpStreamKey) {
	fb.key = key
}

func (fb *fileblk) GetKey() store.UdpStreamKey {
	return fb.key
}
