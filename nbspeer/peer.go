package nbspeer

import (
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/translayer/store"
	"io"
)

type nbspeer struct {
	uid string
	addr address.UdpAddresser
	lastidex int
}

type Peer interface {
	Send(data [] byte,msgtyp int32) error
	SendAnyWay(data []byte,msgtyp int32)
	SendFile(r io.Reader,msgtyp int32) error
	SendFileAnyWay(r io.Reader,msgtyp int32) error

	GetUid() string
}


func (np *nbspeer)getConn() netcommon.UdpConn  {
	conn:=netcommon.GetConnStoreInstance().GetConn(np.uid)

	if conn == nil|| !conn.Status() {
		for {
			ip, port := np.addr.GetAddrS(np.lastidex)
			if ip == "" {
				return nil
			} else {
				conn=netcommon.NewUdpCreateConnection(ip,"",port,0)
				conn.Dial()
				conn.Hello()
				conn.Connect()
				if conn.WaitHello(){
					cs:=netcommon.GetConnStoreInstance()
					cs.Add(np.uid,conn)
					break
				}else{
					np.lastidex ++
					continue
				}
			}
		}

	}
	return conn

}


func (np *nbspeer)Send(data [] byte,msgtyp int32) error  {
	conn:=np.getConn()
	if conn == nil{
		return nbserr.NbsErr{ErrId:nbserr.PEER_CANT_CONNECT,Errmsg:"can't create connection"}
	}
	um:=message.NewUdpMsg(msgtyp,data)
	if d2snd,err:=um.Serialize();err!=nil {
		conn.Send(d2snd)
	}else {
		return err
	}

	ms:=store.NewBlockStoreInstance()
	ms.AddBlock(um)

	//wait reply


	return nil
}

func (np *nbspeer)SendAnyWay(data []byte,msgtyp int32) error{
	conn:=np.getConn()
	if conn == nil{
		return nbserr.NbsErr{ErrId:nbserr.PEER_CANT_CONNECT,Errmsg:"can't create connection"}
	}
	um:=message.NewUdpMsg(msgtyp,data)
	if d2snd,err:=um.Serialize();err!=nil {
		conn.Send(d2snd)
	}else {
		return err
	}

	ms:=store.NewBlockStoreInstance()
	ms.AddBlock(um)

	return nil
}

func (np *nbspeer)SendFile(r io.Reader,msgtyp int32) error  {
	return nil
}

func (np *nbspeer)SendFileAnyWay(r io.Reader,msgtyp int32) error{
	return nil
}

func (np *nbspeer)GetUid() string{
	return np.uid
}