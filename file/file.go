package file

import path2 "path"

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
	GetFileName() string
}

var Save_file_path = ""

func GetSaveFilePath() string {
	if Save_file_path == ""{
		return "/Users/rickey/umfile"
	}else {
		return Save_file_path
	}
}

func SetSaveFilePath(p string)  {
	Save_file_path = p
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

func (fdesc *filedesc)GetFileName() string  {
	path:=fdesc.GetPath()
	if path == "" {
		path = Save_file_path
	}
	filename:=path2.Join(path,fdesc.GetName())

	return filename
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

func NewEmptyFileHead() FileHead  {
	fdesc:=NewFileDesc()
	return NewFileHead(fdesc)
}



