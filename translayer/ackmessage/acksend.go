package ackmessage

func GetAckMessage(sn, pos uint64, otherpos ...uint64) AckMessage {

	ack := &ackmessage{}

	ack.SetSn(sn)
	ack.SetPos(pos)

	arrpos := make([]uint64, 0)
	arrpos = append(arrpos, otherpos...)

	ack.SetResendPos(arrpos)

	return ack

}
