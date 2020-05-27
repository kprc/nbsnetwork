package nbspeer

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
)

type nbspeer struct {
	uid      string
	addr     address.UdpAddresser
	lastidex int
}

type Peer interface {
	SendMessage(data []byte, msgtyp int32) error
	SendMessageAnyWay(data []byte, msgtyp int32)
	SendMessageFile(r io.Reader, msgtyp int32) error
	SendMessageFileAnyWay(r io.Reader, msgtyp int32) error

	GetUid() string
}

func (np *nbspeer) getConn() netcommon.UdpConn {
	conn := netcommon.GetConnStoreInstance().GetConn(np.uid)

	if conn == nil || !conn.Status() {
		for {
			ip, port := np.addr.GetAddrS(np.lastidex)
			if ip == "" {
				return nil
			} else {
				conn = netcommon.NewUdpCreateConnection(ip, "", port, 0)
				conn.Dial()
				conn.Hello()
				conn.Connect()
				if conn.WaitHello() {
					cs := netcommon.GetConnStoreInstance()
					cs.Add(np.uid, conn)
					break
				} else {
					np.lastidex++
					continue
				}
			}
		}

	}
	return conn

}

func (np *nbspeer) SendMessage(data []byte, msgtyp int32) error {
	conn := np.getConn()
	if conn == nil {
		return nbserr.NbsErr{ErrId: nbserr.PEER_CANT_CONNECT, Errmsg: "can't create connection"}
	}
	um := message.NewUdpMsg(msgtyp, data)
	d2snd, err := um.Serialize()
	if err != nil {
		return err
	}

	ms := store.NewBlockStoreInstance()

	ch := make(chan int64, 1)
	chtimeout := make(chan int64, 1)

	nt := tools.GetNbsTickerInstance()

	nt.RegWithTimeOut(&chtimeout, 1200)

	um.SetInformChan(&ch)
	ms.AddBlock(um)

	conn.Send(d2snd)
	defer func() {
		nt.UnReg(&chtimeout)
		ms.DelBlock(um)
	}()

	select {
	case <-ch:
		return nil
	case <-chtimeout:
		return nbserr.NbsErr{ErrId: nbserr.MSG_NO_REPLY, Errmsg: "send message timeout"}
	}
}

func (np *nbspeer) SendMessageAnyWay(data []byte, msgtyp int32) error {
	conn := np.getConn()
	if conn == nil {
		return nbserr.NbsErr{ErrId: nbserr.PEER_CANT_CONNECT, Errmsg: "can't create connection"}
	}
	um := message.NewUdpMsg(msgtyp, data)
	if d2snd, err := um.Serialize(); err != nil {
		conn.Send(d2snd)
	} else {
		return err
	}

	ms := store.NewBlockStoreInstance()
	ms.AddBlock(um)

	return nil
}

func (np *nbspeer) SendMessageFile(r io.Reader, msgtyp int32) error {
	return nil
}

func (np *nbspeer) SendMessageFileAnyWay(r io.Reader, msgtyp int32) error {
	return nil
}

func (np *nbspeer) GetUid() string {
	return np.uid
}
