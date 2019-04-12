package applayer

import (
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/store"
)

type ctrlblk struct {
	rcvblk netcommon.RcvBlock
	um store.UdpMsg
}

type CtrlBlk interface {
	GetRcvBlk() netcommon.RcvBlock
	GetUdpMsg() store.UdpMsg
}

func NewCtrlBlk(rb netcommon.RcvBlock,um store.UdpMsg) CtrlBlk {
	return &ctrlblk{rcvblk:rb,um:um}
}

func (cb *ctrlblk)GetRcvBlk() netcommon.RcvBlock  {
	return cb.rcvblk
}

func (cb *ctrlblk)GetUdpMsg() store.UdpMsg  {
	return cb.um
}

