package send

import (
	"github.com/kprc/nbsnetwork/common/constant"
	"io"
	"net"
)

type udpWriter struct {
	addr *net.UDPAddr
	sock *net.UDPConn
}


type uwReaderSeeker struct {
	data []byte
}

func (uwrs *uwReaderSeeker)Read(p []byte) (n int, err error){
	minlen := len(p)
	if minlen > len(uwrs.data) {
		minlen = len(uwrs.data)
	}
	cplen := copy(p[0:minlen],uwrs.data)
	uwrs.data = uwrs.data[cplen:]
	return cplen,nil
}

func (uwrs *uwReaderSeeker)Seek(offset int64, whence int) (int64, error)  {
	return 0,nil
}


type UdpWriterer interface {
	Send(r io.ReadSeeker) error
	SendBytes(data []byte)
	io.Writer
}

func NewWriter(addr *net.UDPAddr, sock *net.UDPConn) UdpWriterer {
	uw:=&udpWriter{addr:addr,sock:sock}

	return uw
}

func (uw *udpWriter)SendBytes(data []byte)  {
	uwrs := &uwReaderSeeker{data:data}

	uw.Send(uwrs)
}

func (uw *udpWriter)Send(r io.ReadSeeker) error  {
	bd := NewBlockData(r,constant.UDP_MTU)
	bd.SetWriter(uw)

	sd:=NewStoreData(bd)

	bs:=GetInstance()

	bs.AddBlockDataer(bd.GetSerialNo(),sd)

	sd = bs.GetBlockDataer(bd.GetSerialNo())

	sd.GetBlockData().Send()

	return nil
}

func (uw *udpWriter)Write(p []byte) (n int, err error)   {
	return uw.sock.WriteToUDP(p,uw.addr)
}