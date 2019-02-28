package tools

import (
	"github.com/kprc/nbsnetwork/common/list"
	"sync"
	"sync/atomic"
	"time"
)

type nbsticker struct {
	tc *time.Ticker
	llock sync.Mutex
	l list.List
	stop chan int
	wg *sync.WaitGroup
}


type NbsTicker interface {
	Reg(c *chan int64)
	UnReg(c *chan int64)
	Run()
	Stop()
}

type tickV struct {
	c *chan int64
	cnt int
}

var (
	ntinstance NbsTicker
	oncelock sync.Mutex
	gcnt  int64
)


func newNbsTicker() NbsTicker {
	t:=time.NewTicker(time.Millisecond*500)

	l:=list.NewList(func(v1 interface{}, v2 interface{}) int {
		pc1 := v1.(*tickV)
		pc2 := v2.(*tickV)

		if pc1.c == pc2.c {
			return 0
		}else{
			return 1
		}

	})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	return &nbsticker{tc:t,l:l,stop:make(chan int),wg:wg}
}

func GetNbsTickerInstance() NbsTicker  {
	if ntinstance != nil{
		return ntinstance
	}

	oncelock.Lock()
	defer oncelock.Unlock()

	if ntinstance!=nil {
		return ntinstance
	}

	ntinstance = newNbsTicker()

	return ntinstance

}

func (nt *nbsticker)Reg(c *chan int64)  {
	nt.llock.Lock()
	tv := &tickV{c:c}
	nt.l.AddValue(tv)
	nt.llock.Unlock()
}

func (nt *nbsticker)UnReg(c *chan int64)  {
	nt.llock.Lock()
	tv := &tickV{c:c}
	nt.l.DelValue(tv)
	nt.llock.Unlock()
}

func (nt *nbsticker)Stop()  {
	nt.stop <- 0
	if nt.wg !=nil {
		(*nt.wg).Wait()
	}
}

func (nt *nbsticker)delTicker(arr []*tickV)  {
	for _,ptv:=range arr {
		nt.UnReg(ptv.c)
	}
}

func (nt *nbsticker)Run(){
	if nt.wg !=nil{
		defer (*nt.wg).Done()
	}
	for {
		select {
			case <-nt.tc.C:
				arr := make([]*tickV,0)
				nt.llock.Lock()

				nt.l.Traverse(gcnt, func(arg interface{}, data interface{}) {
					cnt := arg.(int64)
					c:= data.(*tickV)
					select {
						case *c.c <- cnt:
							c.cnt = 0
					    default:
					    	c.cnt ++
					    	if c.cnt > 100{
					    		arr = append(arr,c)
							}
					}

				})
				nt.llock.Unlock()
			    atomic.AddInt64(&gcnt,1)
				if len(arr)>0{
					nt.delTicker(arr)
				}
			case <-nt.stop:
				nt.tc.Stop()
				return
		}
	}
}








