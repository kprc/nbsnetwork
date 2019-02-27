package netcommon

type rcvblock struct {
	cp ConnPacket
	conn UdpConn
}


type RcvBlock interface {
	GetConnPacket() ConnPacket
	GetUdpConn() UdpConn
}

func NewRcvBlock(packet ConnPacket,conn UdpConn) RcvBlock  {
	return &rcvblock{cp:packet,conn:conn}
}

func (rb *rcvblock)GetConnPacket() ConnPacket  {
	return rb.cp
}

func (rb *rcvblock)GetUdpConn() UdpConn {
	return rb.conn
}
