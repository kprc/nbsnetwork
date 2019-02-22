package tools

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/list"
	"sync/atomic"
	"sync"
	"time"
)

type nbsticker struct {
	tc *time.Ticker
	llock sync.Mutex
	l list.List
	stop chan int
}


type NbsTicker interface {
	Reg(c *chan int64)
	UnReg(c *chan int64)
	Run(wg *sync.WaitGroup)
	Stop()
}

var (
	ntinstance NbsTicker
	oncelock sync.Mutex
	gcnt  int64
)


func newNbsTicker() NbsTicker {
	t:=time.NewTicker(time.Millisecond*300)

	l:=list.NewList(func(v1 interface{}, v2 interface{}) int {
		pc1 := v1.(*chan int64)
		pc2 := v2.(*chan int64)

		if pc1 == pc2 {
			return 0
		}else{
			return 1
		}

	})

	return &nbsticker{tc:t,l:l,stop:make(chan int)}
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
	nt.l.AddValue(c)
	nt.llock.Unlock()
}

func (nt *nbsticker)UnReg(c *chan int64)  {
	nt.llock.Lock()
	nt.l.DelValue(c)
	nt.llock.Unlock()
}

func (nt *nbsticker)Stop()  {
	nt.stop <- 0
}


func (nt *nbsticker)Run(wg *sync.WaitGroup){
	if wg !=nil{
		defer (*wg).Done()
	}
	for {
		select {
			case <-nt.tc.C:
				nt.llock.Lock()
				nt.l.Traverse(gcnt, func(arg interface{}, data interface{}) {
					cnt := arg.(int64)
					c:= data.(*chan int64)
					select {
						case *c <- cnt:
					    default:
					}

				})
				nt.llock.Unlock()
			    atomic.AddInt64(&gcnt,1)
			case <-nt.stop:
				nt.tc.Stop()
				fmt.Println("Ticker quit...")
				return
		}
	}
}








