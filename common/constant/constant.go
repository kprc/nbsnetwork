package constant

const (
	FILE_DESC_HANDLE   uint32 = 0
	FILE_STREAM_HANDLE uint32 = 1
	FILE_START_SIZE    uint32 = 2
)

const (
	UDP_CONNECTION_TIMEOUT      int = 15000
	UDP_MESSAGE_STORE_TIMEOUT   int = 15000
	UDP_STREAM_STORE_TIMEOUT    int = 15000
	FILE_STORE_TIMEOUT          int = 15000
	RPC_STORE_TIMEOUT           int = 5000
	STREAM_BLOCK_RESEND_TIMEOUT int = 2000
	STREAM_BLOCK_RESEND_TIMES   int = 3
	STREAM_MTU                  int = 544
	STREAM_MAXCACHE             int = 16 * (1 << 10)
)

const (
	OPEN_FILE    int = 1
	CLOSE_FILE   int = 2
	REFRESH_FILE int = 3
)
