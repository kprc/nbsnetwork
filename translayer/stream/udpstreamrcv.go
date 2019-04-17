package stream

import (
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/ackmessage"
	"github.com/kprc/nbsdht/nbserr"
	"reflect"
	"github.com/kprc/nbsnetwork/applayer"
	"github.com/kprc/nbsnetwork/tools"
)

type streamrcv struct {
	udpmsgcache map[uint64]store.UdpMsg
	lastwritepos uint64   //
	finishflag bool			//write finish flag
	toppos uint64		//max pos
	key store.UdpStreamKey
	w io.WriteCloser
}

type StreamRcv interface {
	SetWriter(w io.WriteCloser)
	//read(buf []byte) (int,error)
	addData(um store.UdpMsg) error
	constructResends(ack ackmessage.AckMessage)
	setTopPos(pos uint64)
	write(cb applayer.CtrlBlk) error
	getFinishFlag() bool
}

func (sr *streamrcv)GetKey() store.UdpStreamKey {
	return sr.key
}

func NewStreamRcvWithParam(uid string,sn uint64) StreamRcv{
	sr:=&streamrcv{}
	sr.udpmsgcache = make(map[uint64]store.UdpMsg)
	key:=store.NewUdpStreamKeyWithParam(uid,sn)

	sr.key = key

	return sr
}

func NewStreamRcv(sk store.UdpStreamKey) StreamRcv{
	return NewStreamRcvWithParam(sk.GetUid(),sk.GetSn())
}

func (sr *streamrcv)getFinishFlag() bool  {
	return sr.finishflag
}

func (sr *streamrcv)addData(um store.UdpMsg) error {
	if _,ok:=sr.udpmsgcache[um.GetPos()];ok {
		return nbserr.NbsErr{ErrId:nbserr.ERROR_DEFAULT,Errmsg:"Data exists"}
	}

	sr.udpmsgcache[um.GetPos()] = um

	return nil
}

func (sr *streamrcv)setTopPos(pos uint64)  {
	sr.toppos = pos
}


func (sr *streamrcv)constructResends(ack ackmessage.AckMessage){
	listkeys:=reflect.ValueOf(sr.udpmsgcache).MapKeys()

	if len(listkeys) == 0 {
		return
	}
	var apptyp uint32
	minpos := sr.lastwritepos
	maxpos := minpos
	for _,k:=range listkeys{
		keypos := k.Uint()
		if maxpos < keypos{
			apptyp = sr.udpmsgcache[keypos].GetAppTyp()
			maxpos = keypos
		}
	}

	arrpos:= make([]uint64,0)

	var i uint64
	for i=tools.GetRealPos(minpos); i< tools.GetRealPos(maxpos); i++{
		if _,ok:=sr.udpmsgcache[tools.AssemblePos(i,apptyp)]; !ok{
			arrpos = append(arrpos,i)
		}
	}

	if len(arrpos) > 0 {
		ack.SetResendPos(arrpos)
	}
}


func Recv(rblk netcommon.RcvBlock)  error{

	data:=rblk.GetConnPacket().GetData()
	um:=store.NewUdpMsg(nil,0)

	if err:=um.DeSerialize(data);err!=nil{
		return err
	}

	um.Print()

	sn:=um.GetSn()
	uid:=string(rblk.GetConnPacket().GetUid())

	key:=store.NewUdpStreamKeyWithParam(uid,sn)

	ss:=store.GetStreamStoreInstance()

	fdo := func(arg interface{}, v interface{}) (ret interface{},err error){
		blk:=store.GetStreamBlkAndRefresh(v).(StreamRcv)
		um:=arg.(store.UdpMsg)
		blk.addData(um)
		//write
		if um.GetLastFlag(){
			blk.setTopPos(um.GetPos())
		}
		//construct a ack message
		ack:=ackmessage.GetAckMessage(um.GetSn(),um.GetPos())
		blk.constructResends(ack)

		return ack,nil
	}

	r,_:=ss.FindStreamDo(key,um,fdo)
	if r == nil{
		sr:=NewStreamRcv(key)
		sr.addData(um)
		ss.AddStream(sr)
		r=ackmessage.GetAckMessage(sn,um.GetPos())
	}

	ackdata,_:=r.(ackmessage.AckMessage).Serialize()
	if ackdata !=nil{
		rblk.GetUdpConn().Send(ackdata,store.UDP_ACK)
	}

	fwrite := func(arg interface{},v interface{})(ret interface{},err error) {
		blk:=store.GetStreamBlk(v).(StreamRcv)
		blk.write(arg.(applayer.CtrlBlk))
		if blk.getFinishFlag() {
			return blk,nil
		}
		return
	}

	cb:=applayer.NewCtrlBlk(rblk,um)

	if r,_=ss.FindStreamDo(key,cb,fwrite);r!=nil{
		ss.DelStream(r)
	}

	return nil
}

func (sr *streamrcv)SetWriter(w io.WriteCloser)  {
	sr.w = w
}

//
//func (sr *streamrcv)read(buf []byte) (int,error)  {
//	if sr.finishflag == true{
//		return 0,io.EOF
//	}
//	pos :=sr.lastwritepos
//	defer func() {
//		sr.lastwritepos = pos
//	}()
//	n:=0
//	for{
//
//		if pos > sr.toppos{
//			sr.finishflag = true
//			return n,io.EOF
//		}
//
//		if v,ok:=sr.udpmsgcache[pos];!ok {
//			return n,nil
//		}else{
//			um:=v.(store.UdpMsg)
//			data:=um.GetData()
//			if len(data) ==0 {
//				delete(sr.udpmsgcache,pos)
//				pos ++
//				continue
//			}
//			nc:=copy(buf[n:],data)
//			n+=nc
//			if n >= len(buf){
//				return n,nil
//			}
//			if nc >=len(data){
//				delete(sr.udpmsgcache,pos)
//				pos ++
//				continue
//
//			}else{
//				data = data[nc:]
//				um.SetData(data)
//				return n,nil
//			}
//
//		}
//
//	}
//
//
//	return 0,nil
//}

func (sr *streamrcv)write(cb applayer.CtrlBlk) error  {

	apptyp:=cb.GetUdpMsg().GetAppTyp()

	if sr.w == nil{
		abs:=applayer.GetAppBlockStore()
		if w,err:=abs.Do(apptyp,cb,false);err!=nil{
			return nbserr.NbsErr{ErrId:nbserr.FILE_CANNT_OPEN,Errmsg:"File can't open"}
		}else{
			sr.w = w.(io.WriteCloser)
		}
	}

	if sr.finishflag == true{
		return io.EOF
	}

	pos :=sr.lastwritepos

	if sr.lastwritepos == 0 && apptyp > 0 {
		pos = tools.AssemblePos(pos,apptyp)
	}

	defer func() {
		sr.lastwritepos = pos
	}()

	for{

		if (pos > sr.toppos) && (sr.toppos != 0){
			sr.finishflag = true
			abs:=applayer.GetAppBlockStore()
			abs.Do(apptyp,cb,true)
			return io.EOF
		}

		if v,ok:=sr.udpmsgcache[pos];!ok {
			return nil
		}else{
			um:=v.(store.UdpMsg)
			data:=um.GetData()
			if len(data) ==0 {
				delete(sr.udpmsgcache,pos)
				pos ++
				continue
			}

			nc,err := sr.w.Write(data)

			if nc >=len(data){
				delete(sr.udpmsgcache,pos)
				pos ++

			}else{
				data = data[nc:]
				um.SetData(data)
			}
			if err!=nil{
				return err
			}

		}

	}

	return nil

}
