package store

import "hash/fnv"

type udpstreamkey struct {
	uid string
	sn uint64
}

type UdpStreamKey interface{
	SetUid(uid string)
	GetUid() string
	SetSn(sn uint64)
	GetSn() uint64
	GetKey() UdpStreamKey
	Equals(key UdpStreamKey) bool
}

type StreamKeyInter interface {
	GetKey() UdpStreamKey
}

func NewUdpStreamKey() UdpStreamKey  {
	return &udpstreamkey{}
}

func NewUdpStreamKeyWithParam(uid string, sn uint64) UdpStreamKey  {
	return &udpstreamkey{uid:uid,sn:sn}
}

func (usk *udpstreamkey)SetUid(uid string)  {
	usk.uid = uid
}

func (usk *udpstreamkey)GetUid() string  {
	return usk.uid
}

func (usk *udpstreamkey)SetSn(sn uint64)  {
	usk.sn = sn
}

func (usk *udpstreamkey)GetSn() uint64  {
	return usk.sn
}

func (usk *udpstreamkey)GetKey() UdpStreamKey  {
	return usk
}

func (usk *udpstreamkey)Equals(key UdpStreamKey) bool  {
	if usk.sn == key.GetSn() && usk.uid == key.GetUid() {
		return true
	}

	return false
}




var StreamKeyHash = func(v interface{}) uint {
	sk:=v.(StreamKeyInter)
	s := fnv.New64()
	s.Write([]byte(sk.GetKey().GetUid()))
	h:=s.Sum64()

	h = sk.GetKey().GetSn() << 8 | h

	h = h & 0x7F

	return uint(h)

}

var StreamKeyEquals = func(v1 interface{}, v2 interface{}) int {
	sk1:=v1.(StreamKeyInter)
	sk2:=v2.(StreamKeyInter)

	if sk1.GetKey().GetUid() == sk2.GetKey().GetUid() {
		if sk1.GetKey().GetSn() == sk2.GetKey().GetSn(){
			return 0
		}
	}

	return 1
}

