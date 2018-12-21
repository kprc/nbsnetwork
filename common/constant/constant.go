package constant

const (
	IP_TYPE_IP4 int = 0
	IP_TYPE_IP6 int = 1
	UDP_MTU uint32 = 544     //576-32
	UDP_SEND_TIMEOUT= 3      //3 second
	UDP_RECHECK_TIMEOUT=2    //1 second
	UDP_MAX_CACHE = 64*1024  // 64K


	UDP_SERIAL_MAGIC_NUM = 0x20151031

)

const (
	PING= iota
	ACK
	DATA_TRANSER
)


const (
	MSG_KA = iota
	MSG_STORE
	MSG_FIND_NODE
	MSG_FIND_VALUE
	MSG_PING
)
