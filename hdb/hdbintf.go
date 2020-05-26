package hdb

import "reflect"

type DBCusor struct {
	keys []reflect.Value
	cursor int
	hdb HistoryDBIntf
}

type HistoryDBIntf interface {
	Load() HistoryDBIntf
	Insert(key string,value string) error
	Delete(key string)
	Find(key string,start int, topn int) ([]*HDBV,error)
	FindBlock(key string) (*FileHDBV,error)
	Save()
	DBIterator() *DBCusor
}

func (dbc *DBCusor)Next() (k string,v []*HDBV)  {
	if dbc.cursor >= len(dbc.keys){
		return
	}

	k = dbc.keys[dbc.cursor].Interface().(string)
	v,_=dbc.hdb.Find(k,0,0)

	dbc.cursor ++

	return
}

