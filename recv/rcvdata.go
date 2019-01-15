package recv

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/flowkey"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/netcommon"
	"io"
	"reflect"
	"sort"
	"sync"
	"time"
)

var rcvdataerr = nbserr.NbsErr{ErrId:nbserr.UDP_RCV_DEFAULT_ERR,Errmsg:"Please initial the node "}
var rcvwriteioerr = nbserr.NbsErr{ErrId:nbserr.UDP_RCV_WRITER_IO_ERR,Errmsg:"write error"}

type rcvData struct {
	key *flowkey.FlowKey
	uw netcommon.UdpReaderWriterer   //for reply
	w io.WriteSeeker
	lastAccessTime int64
	rwlock sync.RWMutex
	rcvData map[uint32]packet.UdpPacketDataer     //cache receive data for write to w
	lastRcvId uint32
	lastWriteId uint32
	totalCnt uint32
	rwlockResend sync.RWMutex
	needResend map[uint32]struct{}    //for caculate what sn is need to resend
	finished bool
	canDelete bool
}


type RcvDataer interface {
	Write(dataer packet.UdpPacketDataer) (packet.UdpResulter,error)
	Finish() bool
	CanDelete() bool
	TimeOut() bool
}

func (rd *rcvData)CanDelete() bool  {
	return rd.canDelete
}

func (rd *rcvData)TimeOut()  bool{

	return false
}

func NewRcvDataer(stationId string,sn uint64,w io.WriteSeeker,uw netcommon.UdpReaderWriterer) RcvDataer{
	fk := flowkey.NewFlowKey(stationId,sn)
	rd := &rcvData{key:fk,w:w,uw:uw}
	rd.lastAccessTime = time.Now().Unix()
	rd.rcvData = make(map[uint32]packet.UdpPacketDataer)
	rd.needResend = make(map[uint32]struct{})

	return rd
}

func (rd *rcvData)write() error  {

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
			if upr,ok := rd.rcvData[id];ok{
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
	}

	if rd.totalCnt>0 && rd.lastWriteId == rd.totalCnt {
		rd.finished = true
	}

	return nil
}

func max(a,b uint32) uint32 {
	if a>b{
		return a
	}else {
		return b
	}
}

func (rd *rcvData)Write(dataer packet.UdpPacketDataer) (packet.UdpResulter,error)  {
	if rd.serialNo==0 || rd.w==nil || dataer == nil {
		return nil,rcvdataerr
	}

	if dataer.IsFinished() {
		rd.canDelete = true
		return nil,nil
	}

	var ack packet.UdpResulter

	rd.lastAccessTime = time.Now().Unix()

	if dataer.GetTotalCnt() > 0 {
		rd.totalCnt = max(rd.totalCnt,dataer.GetTotalCnt())
	}
	rd.enqueen(dataer)
	ack=rd.fillNeedResend(dataer)

	err := rd.write()

	if rd.finished {
		ack.Finished()
	}

	return ack,err
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

func (rd *rcvData)Finish() bool {
	return rd.finished
}