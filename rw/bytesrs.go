package rw

import (
	"io"
	"fmt"
)

type uwReaderSeeker struct {
	data []byte
}

type UdpBytesReaderSeeker interface {
	io.ReadSeeker
	PrintAll()
	GetBytes() []byte
}


func NewReadSeeker(data []byte) UdpBytesReaderSeeker {
	return &uwReaderSeeker{data:data}
}

func (uwrs *uwReaderSeeker)Read(p []byte) (n int, err error){
	minlen := len(p)
	if minlen > len(uwrs.data) {
		minlen = len(uwrs.data)
	}

	if minlen == 0 {
		 return 0,io.EOF
	}
	cplen := copy(p[0:minlen],uwrs.data)
	uwrs.data = uwrs.data[cplen:]
	return cplen,nil
}

func (uwrs *uwReaderSeeker)Seek(offset int64, whence int) (int64, error)  {
	return 0,nil
}

func (uwrs *uwReaderSeeker) PrintAll() {
	fmt.Println(string(uwrs.data))
}

func (uwrs *uwReaderSeeker)GetBytes() []byte  {
	return uwrs.data
}
