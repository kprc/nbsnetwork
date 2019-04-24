package bus

import (
	"sync"
	"github.com/kprc/nbsdht/nbserr"
)

type busstore struct {
	bus chan interface{}
	quit chan int
	wg *sync.WaitGroup
}

type BusStore interface {
	AddBlock(v interface{}) error
	Run()
	Stop()
}

var (
	busstoreInstance BusStore
	busstoreInstLock sync.Mutex
)


var (
	busoverflow = nbserr.NbsErr{ErrId:nbserr.BUS_OVERFLOW,Errmsg:"Bus Overflow, too many data to send"}
)

func newBusStore() BusStore {
	bs:=&busstore{}
	bs.bus = make(chan interface{},2048)
	bs.quit= make(chan int, 1)
	bs.wg = &sync.WaitGroup{}

	return bs
}

func GetBusStoreInstance() BusStore  {
	if busstoreInstance == nil{
		busstoreInstLock.Lock()
		defer busstoreInstLock.Unlock()

		if busstoreInstance == nil{
			busstoreInstance = newBusStore()
		}

	}


	return busstoreInstance
}

func (bs *busstore)AddBlock(v interface{}) error {
	select {
		case bs.bus <- v:
		default:
			return busoverflow
	}
	return nil
}

func (bs *busstore)Run()  {
	if bs.wg!=nil{
		bs.wg.Add(1)
	}
	defer func() {
		if bs.wg != nil{
			bs.wg.Done()
		}
	}()

	for{
		select {
		case v:=<-bs.bus:
			SendBusBlock(v)
			case <-bs.quit:
				return
		}
	}
}

func (bs *busstore)Stop()  {
	bs.quit <- 0
	if bs.wg != nil{
		bs.wg.Wait()
	}
}





