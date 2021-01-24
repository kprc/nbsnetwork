package httputil

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

func Post(url string, jsonstr string, blog bool) (jsonret string, code int, err error) {
	if blog {
		log.Println(url)
		log.Println(jsonstr)
	}

	bjson := []byte(jsonstr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bjson))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errresp := client.Do(req)

	if errresp != nil {
		return "", 0, errresp
	}

	defer resp.Body.Close()

	body, errbody := ioutil.ReadAll(resp.Body)
	if errbody != nil {
		return "", 0, errbody
	}

	return string(body), resp.StatusCode, nil
}

func GetRemoteAddr(remoteaddr string) (ip string, port string) {
	if remoteaddr == "" {
		return
	}

	ra := strings.Split(remoteaddr, ":")

	if len(ra) != 2 {
		return
	}

	ip = ra[0]
	port = ra[1]

	return
}

type HttpPost struct {
	Protect func(fd int32) bool
	Blog bool
	DialTimeout int
	ConnTimeout int
}

func NewHttpPost(protect func(fd int32)bool, blog bool, dialTimeout, ConnTimeout int) *HttpPost {
	return &HttpPost{
		Protect: protect,
		Blog: blog,
		DialTimeout: dialTimeout,
		ConnTimeout: dialTimeout,
	}
}

func (hp *HttpPost)ProtectPost(url string, jsonstr string) (jsonret string, code int, err error)   {
	if hp.Blog {
		log.Println(url)
		log.Println(jsonstr)
	}

	bjson := []byte(jsonstr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bjson))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	var transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error){
			d := &net.Dialer{
				Timeout: time.Second * time.Duration(hp.DialTimeout),
				Control: func(network, address string, c syscall.RawConn) error {
					if hp.Protect != nil {
						p:= func(fd uintptr) {
							hp.Protect(int32(fd))
						}
						return c.Control(p)
					}
					return nil
				},
			}

			conn, err := d.Dial("udp", addr)
			if err != nil {
				return nil, err
			}

			if hp.ConnTimeout > 0{
				conn.SetDeadline(time.Now().Add(time.Duration(hp.ConnTimeout)*time.Second))
			}

			return conn,nil
		},
	}

	client := &http.Client{
		Transport: transport,
	}
	resp, errresp := client.Do(req)

	if errresp != nil {
		return "", 0, errresp
	}

	defer resp.Body.Close()

	body, errbody := ioutil.ReadAll(resp.Body)
	if errbody != nil {
		return "", 0, errbody
	}

	return string(body), resp.StatusCode, nil
}