package privateip

import (
	"bytes"
	"net"
	"strings"
	"sync"
)

//address from https://en.wikipedia.org/wiki/Reserved_IP_addresses
var NetWorkIPs = []string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.88.99.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"255.255.255.255/32",
}

type NetworkIP struct {
	netIP *net.IP
	mask  *net.IPMask
}

var gNetworkArr []*NetworkIP
var gNetworkArrLock sync.Mutex

func init() {
	if gNetworkArr != nil {
		return
	}

	gNetworkArrLock.Lock()
	defer gNetworkArrLock.Unlock()
	if gNetworkArr != nil {
		return
	}

	gNetworkArr = make([]*NetworkIP, 0)

	for _, sip := range NetWorkIPs {
		r := convertP(sip)
		if r == nil {
			continue
		}

		gNetworkArr = append(gNetworkArr, r)
	}
}

func convertP(sip string) *NetworkIP {
	siparr := strings.Split(sip, "/")

	if len(siparr) != 2 {
		return nil
	}
	ip, ipnet, err := net.ParseCIDR(sip)
	if err != nil {
		return nil
	}

	r := &NetworkIP{netIP: &ip, mask: &ipnet.Mask}

	return r
}

func IsPrivateIP(ip net.IP) bool {
	if nil == ip {
		return false
	}
	for _, ni := range gNetworkArr {
		netip := ip.Mask(*ni.mask)
		if bytes.Compare([]byte(netip.To4()), []byte((*ni.netIP).To4())) == 0 {
			return true
		}
	}
	return false
}

func IsPrivateIPStr(ips string) bool {

	ip := net.ParseIP(ips)

	return IsPrivateIP(ip)
}
