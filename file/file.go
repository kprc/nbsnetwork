package file

type filedesc struct {
	path string
	name string
	size uint64
}

type FileDesc interface {
	GetPath() string
	SetPath(path string)
	GetName() string
	SetName(name string)
	GetSize() uint64
	SetSize(size uint64)
}


func (fdesc *filedesc)GetPath() string  {
	return fdesc.path
}

func (fdesc *filedesc)SetPath(path string)  {
	fdesc.path = path
}

func (fdesc *filedesc)GetName() string {
	return fdesc.name
}

func (fdesc *filedesc)SetName(name string)  {
	fdesc.name = name
}

func (fdesc *filedesc)GetSize() uint64  {
	return fdesc.size
}

func (fdesc *filedesc)SetSize(size uint64)  {
	fdesc.size = size
}

func NewFileDesc() FileDesc {
	return &filedesc{}
}


type filehead struct {
	FileDesc
	strhash string
}

type FileHead interface {
	FileDesc
	GetStrHash() string
	SetStrHash(sh string)
}

func (fh *filehead)GetStrHash() string  {
	return fh.strhash
}

func (fh *filehead)SetStrHash(sh string)  {
	fh.strhash = sh
}

func NewFileHead(desc FileDesc) FileHead {
	return &filehead{desc,""}
}




