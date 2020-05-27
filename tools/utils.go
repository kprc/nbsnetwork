package tools

import (
	"github.com/pkg/errors"
	"net"
	"strconv"
	"strings"
	"time"
)

func ResizeHash(hash uint) uint {
	cnt := 0
	for {
		hash = hash >> 1
		if hash > 0 {
			cnt++
		} else {
			break
		}
	}
	return 1 << uint(cnt)
}

func GetNowMsTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetRealPos(pos uint64) uint64 {
	return (pos & 0xFFFFFFFF)
}

func AssemblePos(pos uint64, typ uint32) uint64 {
	var typ1 uint64

	typ1 = uint64(typ)

	typ1 = typ1 << 32

	pos = pos & 0xFFFFFFFF

	return typ1 | pos

}

func GetTypFromPos(pos uint64) uint32 {
	typ := (pos >> 32) & 0xFFFFFFFF

	return uint32(typ)
}

func CheckPortUsed(iptyp, ipaddr string, port uint16) bool {
	if strings.Contains(strings.ToLower(iptyp), "udp") {
		netaddr := &net.UDPAddr{IP: net.ParseIP(ipaddr), Port: int(port)}
		if c, err := net.ListenUDP(iptyp, netaddr); err != nil {
			return true
		} else {
			c.Close()
			return false
		}
	} else {
		netaddr := &net.TCPAddr{IP: net.ParseIP(ipaddr), Port: int(port)}
		if c, err := net.ListenTCP(iptyp, netaddr); err != nil {
			return true
		} else {
			c.Close()
			return false
		}
	}
}

func GetIPPort(addr string) (ip string, port int, err error) {
	arraddr := strings.Split(addr, ":")
	if len(arraddr) != 2 {
		return "", 0, errors.New("address error")
	}

	ip = arraddr[0]
	port, err = strconv.Atoi(arraddr[1])
	if err != nil {
		return "", 0, err
	}
	if port < 1024 || port > 65535 {
		return "", 0, errors.New("port error")
	}

	if _, err = net.ResolveIPAddr("ip4", ip); err != nil {
		return "", 0, err
	}

	return ip, port, nil
}
