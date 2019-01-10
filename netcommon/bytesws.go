package netcommon

import (
	"fmt"
	"io"
)

type uwWriterSeeker struct {
	data []byte
}

type UdpBytesWriterSeeker interface {
	io.WriteSeeker
	PrintAll()
	GetBytes() []byte
}


func NewWriteSeeker(param interface{})UdpBytesWriterSeeker {
	return &uwWriterSeeker{data:make([]byte,0)}
}

func (uwrs *uwWriterSeeker)Write(p []byte) (n int, err error){
	uwrs.data = append(uwrs.data,p...)

	return len(p),nil
}

func (uwrs *uwWriterSeeker)Seek(offset int64, whence int) (int64, error)  {
	return 0,nil
}

func (uwrs *uwWriterSeeker)PrintAll()  {
	fmt.Println(string(uwrs.data))
}

func (uwrs *uwWriterSeeker)GetBytes() []byte  {
	return uwrs.data
}
