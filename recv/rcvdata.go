package recv

import (
	"sync"
	"github.com/kprc/nbsnetwork/common/packet"
	"io"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsdht/nbserr"
	"time"
	"sort"
	"reflect"
)

var rcvdataerr = nbserr.NbsErr{ErrId:nbserr.UDP_RCV_DEFAULT_ERR,Errmsg:"Please initial the node "}
var rcvwriteioerr = nbserr.NbsErr{ErrId:nbserr.UDP_RCV_WRITER_IO_ERR,Errmsg:"write error"}

type rcvData struct {
	serialNo uint64
	lastAccess int64
	timeout uint32
	rwlock sync.RWMutex
	rcvData map[uint32]packet.UdpPacketDataer
	w io.WriteSeeker
	lastRcvId uint32
	lastWriteId uint32
	rwlockResend sync.RWMutex
	needResend map[uint32]struct{}
	finish bool
}


type RcvDataer interface {
	Write(dataer packet.UdpPacketDataer) (packet.UdpResulter,error)
}

func NewRcvDataer(sn uint64,w io.WriteSeeker) RcvDataer{
	rd := &rcvData{serialNo:sn,w:w}
	rd.timeout = constant.UDP_RECHECK_TIMEOUT
	rd.rcvData = make(map[uint32]packet.UdpPacketDataer)
	rd.needResend = make(map[uint32]struct{})

	return rd
}

func (rd *rcvData)write(dataer packet.UdpPacketDataer) error  {

	rd.rwlock.Lock()
	defer rd.rwlock.Unlock()
	keys:=reflect.ValueOf(rd.rcvData).MapKeys()

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Uint() < keys[j].Uint() {
			return true
		}else {
			return false
		}
	})

	rd.w.Seek(0,io.SeekEnd)

	for _,k:=range keys {
		id := uint32(k.Uint())
		if rd.lastWriteId >= id {
			delete(rd.rcvData,id)
		}else if id == rd.lastWriteId + 1 {
			upr := rd.rcvData[id]

			if nw,err := rd.w.Write(upr.GetData());err!=nil{
				if nw>0 {
					needwrite := upr.GetData()
					upr.SetData(needwrite[nw:])
					upr.SetLength(int32(len(upr.GetData())))
				}
				return rcvwriteioerr
			}
			delete(rd.rcvData,id)
			rd.lastWriteId = id

		}
	}




	return nil
}

func (rd *rcvData)Write(dataer packet.UdpPacketDataer) (packet.UdpResulter,error)  {
	if rd.serialNo==0 || rd.w==nil {
		return nil,rcvdataerr
	}


	rd.lastAccess = time.Now().Unix()
	rd.enqueen(dataer)

    ack:=rd.fillNeedResend(dataer)

	reterr := rd.write(dataer)

	return ack,reterr
}



func (rd *rcvData)fillNeedResend(dataer packet.UdpPacketDataer) packet.UdpResulter{

	rd.rwlockResend.Lock()
	defer rd.rwlockResend.Unlock()

	ur := packet.NewUdpResult(dataer.GetSerialNo())

	rd.needResend[dataer.GetPos()] = struct{}{}

	keys := reflect.ValueOf(rd.needResend).MapKeys()

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Uint() < keys[j].Uint() {
			return true
		}else {
			return false
		}
	})

	var max uint32
	if len(keys) > 0 {
		max = uint32(keys[len(keys)-1].Uint())
	}

	for _,k :=range keys{
		id:=uint32(k.Uint())
		if rd.lastRcvId >= id {
			delete(rd.needResend,id)
		}else {
			if rd.lastRcvId +1 == id {
				rd.lastRcvId ++
				delete(rd.needResend,id)
			}
		}
	}

	id:=rd.lastRcvId+1
	for{
		if id<=max {
			if _,ok:=rd.needResend[id];!ok {
				ur.AppendResend(id)
			}
		}else {
			break
		}
		id++
	}

	ur.SetRcved(dataer.GetPos())

	return ur
}


func (rd *rcvData)enqueen(upd packet.UdpPacketDataer){
	id := upd.GetPos()
	rd.rwlock.Lock()
	defer rd.rwlock.Unlock()
	if _,ok:=rd.rcvData[id];!ok{
		rd.rcvData[id] = upd
	}
}