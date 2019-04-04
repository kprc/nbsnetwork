package stream

import (
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/pb/udpmessage"
	"github.com/gogo/protobuf/proto"
)

type streamrcv struct {
	udpmsgcache map[uint64]store.UdpMsg
	lastwritepos uint64
	finishflag bool
	toppos uint64
	key store.UdpStreamKey
	w io.Writer
}

type StreamRcv interface {
	Recv(rblk netcommon.RcvBlock) error
	SetWriter(w io.Writer)
	Read(buf []byte) error
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


func (sr *streamrcv)Recv(rblk netcommon.RcvBlock)  error{
	data:=rblk.GetConnPacket().GetData()
	um:=&udpmessage.Udpmsg{}

	if err:=proto.Unmarshal(data,um);err!=nil{
		return err
	}

	sn:=um.GetSn()
	uid:=string(rblk.GetConnPacket().GetUid())

	key:=store.NewUdpStreamKeyWithParam(uid,sn)

	ss:=store.GetStreamStoreInstance()

	fdo := func(arg interface{}, v interface{}) (ret interface{},err error){
		return v,nil
	}

	if r,_:=ss.FindStreamDo(key,nil,fdo); r==nil{
			
	}

	return nil
}

func (sr *streamrcv)SetWriter(w io.Writer)  {
	sr.w = w
}

func (sr *streamrcv)Read(buf []byte) error  {
	return nil
}
