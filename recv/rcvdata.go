package recv

import (
	"sync"
	"github.com/kprc/nbsnetwork/common/packet"
	"io"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsdht/nbserr"
)

var rcvdataerr = nbserr.NbsErr{ErrId:nbserr.UDP_RCV_DEFAULT_ERR,Errmsg:"Please initial the node "}

type rcvData struct {
	serialNo uint64
	timeout uint32
	rwlock sync.RWMutex
	rcvData map[uint32]packet.UdpPacketDataer
	w io.WriteSeeker
	lastRcvId uint32
	lastWriteId uint32
	needResend map[uint32]bool
	finish bool
}




type RcvDataer interface {
	Write(dataer packet.UdpPacketDataer) error
}

func NewRcvDataer(sn uint64,w io.WriteSeeker) RcvDataer{
	rd := &rcvData{serialNo:sn,w:w}
	rd.timeout = constant.UDP_RECHECK_TIMEOUT
	rd.rcvData = make(map[uint32]packet.UdpPacketDataer)
	rd.needResend = make(map[uint32]bool)

	return rd
}

func (rd *rcvData)Write(dataer packet.UdpPacketDataer) error  {
	if rd.serialNo==0 || rd.w==nil {
		return rcvdataerr
	}



}


func (rd *rcvData)enqueen(id uint32,upd packet.UdpPacketDataer){
	rd.rwlock.Lock()
	defer rd.rwlock.Unlock()
	if _,ok:=rd.rcvData[id];!ok{
		rd.rcvData[id] = upd
	}
}