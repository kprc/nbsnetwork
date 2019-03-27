package stream


type udpstream struct {

}


type UdpStream interface {
	ReadAsync(rcvblk []byte) error
	Read(rcvblk []byte, timeout int) error
}






