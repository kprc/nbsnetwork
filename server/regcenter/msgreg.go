package regcenter

import "sync"

type msgCenter struct {
	rwlock sync.RWMutex
	coor map[int]MsgHandler
}


type MsgCenter interface {

}

