package nbsnetwork

const (
	IP_TYPE_IP4 int = 0
	IP_TYPE_IP6 int = 1
	UDP_MTU uint32 = 576

	UDP_SERIAL_MAGIC_NUM = 0x20151030


)

const (
	PING= iota
	ACK
	DATA_TRANSER
)

