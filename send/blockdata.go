package send

import (
	"fmt"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"github.com/kprc/nbsnetwork/pb/message"
	"io"
	"sync"
	"sync/atomic"
	"time"
	"github.com/gogo/protobuf/proto"
)

var blocksnderr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_DEFAULT_ERR,Errmsg:"Send error"}
var blocksndmtuerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_MTU_ERR,Errmsg:"mtu is 0"}
var blocksndreaderioerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_READER_IO_ERR,Errmsg:"Reader io error"}
var blocksndwriterioerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_WRITER_IO_ERR,Errmsg:"Writer io error"}
var blocksndtimeout = nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"send timeout"}

type BlockData struct {
	lastSendTime int64   //for timeout
	r io.ReadSeeker      //for read content
	w io.Writer			 //for snd io
	serialNo uint64      //for search the block
	mtu   uint32
	totalreadlen uint32   //for seek
	noacklen int32
	maxcache int32       //for udp local stack
	curNum uint32
	totalSndCnt uint32
	dataType uint16
	transinfo []byte
	chResult chan interface{}
	cmd chan int
	waitAck chan int
	rwlock sync.RWMutex
	sndData map[uint32]packet.UdpPacketDataer
	finished bool
}

type BlockDataer interface {
	Send() error
	SetWriter(w io.Writer)
	GetSerialNo() uint64
	PushResult(result interface{})
	SetDataTyp(typ uint16)
	GetDataTyp() uint16
	SetTransInfo(info []byte)
	GetTransInfo() []byte
	SetTransInfoCommon(stationId string,msgid int32) error
	GetTransInfoCommon() (stationId string, msgid int32,err error)
	SetTransInfoHead(head []byte) error
	GetTransInfoHead() (head []byte,err error)
	SetTransInfoOrigin(stationId string,msgid int32,head []byte) error
	GetTransInfoOrigin() (stationId string,msgid int32,head []byte,err error)
	Finished()
	IsFinished() bool
	TimeOut()
	Destroy()
	SendAll()
}

var gSerialNo uint64 = constant.UDP_SERIAL_MAGIC_NUM

func (uh *BlockData)nextSerialNo() {
	uh.serialNo = atomic.AddUint64(&gSerialNo,1)
}

func NewBlockData(r io.ReadSeeker,mtu uint32) BlockDataer {
	uh := &BlockData{r:r,mtu:mtu}
	uh.nextSerialNo()
	uh.lastSendTime = time.Now().Unix()
	uh.mtu = constant.UDP_MTU
	uh.maxcache = constant.UDP_MAX_CACHE
	uh.chResult = make(chan interface{},1024)
	uh.cmd = make(chan int,0)
	uh.waitAck = make(chan int, 0)

	uh.sndData = make(map[uint32]packet.UdpPacketDataer)
	return uh
}

func (bd *BlockData)GetSerialNo() uint64  {
	return bd.serialNo
}

func (bd *BlockData)SetWriter(w io.Writer)  {
	bd.w = w
}

func (bd *BlockData)send(round uint32) (int,error){
	buf := make([]byte,bd.mtu)

	n,err := bd.r.Read(buf)

	if n > 0 {
		fmt.Println("read:===>",string(buf[:n]))
		atomic.AddInt32(&bd.noacklen,int32(n))
		upr := packet.NewUdpPacketData()
		upr.SetTyp(bd.dataType)
		upr.SetSerialNo(bd.serialNo)
		upr.SetData(buf[:n])
		upr.SetTransInfo(bd.transinfo)

		upr.SetTotalCnt(0)
		upr.SetPos(round)
		upr.SetLength(int32(n))
		atomic.AddUint32(&bd.totalreadlen,uint32(n))

		if err==io.EOF {
			upr.SetTotalCnt(round)
		}
		bupr,_ := upr.Serialize()
		fmt.Println(len(bupr))
		nw,err := bd.w.Write(bupr)

		if len(bupr)!=nw ||  err!=nil{
			//need resend
			bd.enqueue(round,upr)
			return 0,blocksndwriterioerr
		}
		atomic.AddUint32(&bd.totalSndCnt,1)
		bd.enqueue(round,upr)
	}

	if err == io.EOF {
		return 1,nil
	}else if err!=nil {
		return 0,blocksndreaderioerr
	}

	return 0,nil
}

func (bd *BlockData)continuesend() (uint32,error) {

	var round uint32 = 1
	bd.rwlock.RLock()
	defer bd.rwlock.RUnlock()

	for _,upr := range bd.sndData{
		bupr,_ := upr.Serialize()
		nw,err := bd.w.Write(bupr)
		if round < upr.GetPos() {
			round = upr.GetPos()
		}

		if len(bupr)!=nw ||  err!=nil{
			//need resend

			return 0,blocksndwriterioerr
		}

		upr.IncTryCnt()
		atomic.AddUint32(&bd.totalSndCnt,1)
	}

	return round,nil
}

func (bd *BlockData)SendAll()  {
	go bd.doACK()
	bd.Send()
}

func (bd *BlockData)Send() error {

	ret := 0

	if bd.r == nil || bd.w == nil{
		return blocksnderr
	}

	if bd.mtu == 0 {
		return blocksndmtuerr
	}

	round,err := bd.continuesend()
	if err != nil{
		return err
	}

	if bd.totalreadlen > 0 {
		if _, err := bd.r.Seek(int64(bd.totalreadlen), io.SeekStart); err != nil {
			return blocksndreaderioerr
		}
	}

	for {
		bd.lastSendTime = time.Now().Unix()
		if ret ==0 && atomic.LoadInt32(&bd.noacklen) < int32(bd.maxcache){
			if r,err := bd.send(round); err==nil{
				round++
				ret = r
			}else if err!=nil{
				return err
			}
		}

		if ret == 1{
			return nil
		}

		wa:=<-bd.waitAck

		if  wa == 1{
			return nil
		}

	}

}

func (bd *BlockData)sendFinish() error {
	upd:=packet.NewUdpPacketData()
	upd.SetFinishACK()
	upd.SetSerialNo(bd.GetSerialNo())

	mc:=regcenter.GetMsgCenterInstance()
	msgid,_,hi:=mc.GetMsgId(bd.GetTransInfo())


	mh:=message.MsgHead{LocalStationId:[]byte(nbsid.GetLocalId().String()),MsgId:msgid}
	mh.Headinfo = hi

	bmh,err := proto.Marshal(&mh)
	if err!=nil {
		return err
	}
	upd.SetTransInfo(bmh)

	bupd,err := upd.Serialize()

	if err!=nil {
		return err
	}

	bd.w.Write(bupd)

	bd.Finished()

	return nil
}

func (bd *BlockData)doresult(result interface{}) error {

	switch v:=result.(type) {
	case []byte:
		r := packet.NewUdpResult(bd.GetSerialNo())
		if err := r.DeSerialize(v); err!=nil{
			return err
		}

		if r.IsFinished() {
			bd.sendFinish()
			return nil
		}
		bd.rwlock.RLock()
		resend := r.GetReSend()
		for _, id := range resend {
			if upr, ok := bd.sndData[id]; ok {
				bupr, _ := upr.Serialize()
				nw, err := bd.w.Write(bupr)

				if len(bupr) != nw || err != nil {
					//need resend
					continue
				}
				upr.SetTryCnt(upr.GetTryCnt() + 1)
				atomic.AddUint32(&bd.totalSndCnt, 1)
			}
		}
		bd.rwlock.RUnlock()
		bd.rwlock.Lock()
		if r.GetRcved() > 0 {
			if upr, ok := bd.sndData[r.GetRcved()]; ok {
				atomic.AddInt32(&bd.noacklen, upr.GetLength()*(-1))
			}
			delete(bd.sndData, r.GetRcved())
		}
		bd.rwlock.Unlock()
		bd.waitAck <- 0    // continue send
	}

	return nil
}

func (bd *BlockData)PushResult(result interface{})  {
	bd.chResult <- result
}

func (bd *BlockData)doACK() {
	for {
		select {
		case result := <-bd.chResult:
			bd.doresult(result)
		case <-bd.cmd:
			bd.finished = true
			return
		}
	}
}

func (bd *BlockData)TimeOut()  {
	if time.Now().Unix() - bd.lastSendTime > 3 {
		bd.Notify()
	}

}


func (bd *BlockData)enqueue(pos uint32, data packet.UdpPacketDataer)  {
	bd.rwlock.Lock()
	defer bd.rwlock.Unlock()
	bd.sndData[pos] = data
}


func (bd *BlockData)SetTransInfo(info []byte){
	bd.transinfo = info

}
func (bd *BlockData)GetTransInfo() []byte{
	return bd.transinfo
}
func (bd *BlockData)SetTransInfoCommon(stationId string,msgid int32) error{
	mh:=message.MsgHead{LocalStationId:[]byte(stationId),MsgId:msgid}

	bmh,err := proto.Marshal(&mh)
	if err!=nil {
		return err
	}

	bd.transinfo = bmh

	return nil
}
func (bd *BlockData)GetTransInfoCommon() (stationId string, msgid int32,err error){
	mh:=message.MsgHead{}

	err = proto.Unmarshal(bd.transinfo,&mh)

	if err!=nil {
		return "",0,err
	}

	return string(mh.LocalStationId),mh.MsgId,err

}
func (bd *BlockData)SetTransInfoHead(head []byte) error{
	if bd.transinfo == nil {
		return nbserr.NbsErr{Errmsg:"Transfer info is nil",ErrId:nbserr.UDP_TRANSINFO_ERR}
	}
	mh:=message.MsgHead{}

	err := proto.Unmarshal(bd.transinfo,&mh)
	if err!=nil {
		return err
	}
	mh.Headinfo = head

	bmh,err1:=proto.Marshal(&mh)
	if err1 != nil {
		return err1
	}

	bd.transinfo = bmh

	return nil
}
func (bd *BlockData)GetTransInfoHead() (head []byte,err error){
	if bd.transinfo == nil {
		return nil,nbserr.NbsErr{Errmsg:"Transfer info is nil",ErrId:nbserr.UDP_TRANSINFO_ERR}
	}

	mh:=message.MsgHead{}

	err = proto.Unmarshal(bd.transinfo,&mh)
	if err!=nil {
		return nil,err
	}

	head = mh.Headinfo

	return
}


func (bd *BlockData)SetTransInfoOrigin(stationId string,msgid int32,head []byte) error {
	mh:=message.MsgHead{MsgId:msgid}


	mh.Headinfo = make([]byte,len(head))
	copy(mh.Headinfo,head)
	mh.LocalStationId = make([]byte,len([]byte(stationId)))
	copy(mh.LocalStationId,[]byte(stationId))

	bmh,err:=proto.Marshal(&mh)
	if err!=nil {
		return err
	}

	bd.transinfo = bmh

	return nil
}
func (bd *BlockData)GetTransInfoOrigin() (stationId string,msgid int32,head []byte,err error){
	if bd.transinfo == nil {
		return "",0,nil,nbserr.NbsErr{Errmsg:"Transfer info is nil",ErrId:nbserr.UDP_TRANSINFO_ERR}
	}
	mh:=message.MsgHead{}

	err = proto.Unmarshal(bd.transinfo,&mh)

	return string(mh.LocalStationId),mh.MsgId,mh.Headinfo,err

}

func (bd *BlockData)SetDataTyp(typ uint16){
	bd.dataType = typ
}

func (bd *BlockData)GetDataTyp() uint16  {
	return bd.dataType
}

func (bd *BlockData)Finished()  {
	bd.cmd <- 1
}

func (bd *BlockData)IsFinished() bool {
	return bd.finished
}

func (bd *BlockData)Notify()  {
	bd.waitAck <- 1
	bd.cmd <- 1
}

func (bd *BlockData)Destroy()  {
	close(bd.waitAck)
	close(bd.chResult)
	close(bd.cmd)
}

