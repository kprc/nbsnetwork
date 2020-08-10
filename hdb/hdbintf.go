package hdb

import (
	"reflect"
)

type DBCusor struct {
	keys   []reflect.Value
	cursor int
	hdb    HistoryDBIntf
}

type HistoryDBIntf interface {
	Load() HistoryDBIntf
	Insert(key string, value string) (int, error)
	Delete(key string)
	FindMem(key string, start int, n int) ([]*HDBV, error)
	FindLatest(key string) (*HDBV, error)
	Find(key string, start, n int) ([]*HDBV, error)
	FindBlock(key string) (*FileHDBV, error)
	//TrimHeadCount(key string, n int)
	//TrimHeadPosition(key string, pos int)
	Save()
	DBIterator() *DBCusor
}

func (dbc *DBCusor) Next() (k string, v []*HDBV) {
	if dbc.cursor >= len(dbc.keys) {
		return
	}

	//var err error

	k = dbc.keys[dbc.cursor].Interface().(string)
	v, _ = dbc.hdb.FindMem(k, 0, 0)
	//if err!=nil{
	//	fmt.Println(err)
	//}

	dbc.cursor++

	return
}
