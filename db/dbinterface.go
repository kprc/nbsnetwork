package db

import "reflect"

type DBCusor struct {
	keys   []reflect.Value
	cursor int
	db     NbsDbInter
}

type NbsDbInter interface {
	Load() NbsDbInter
	Insert(key string, value string) error
	Delete(key string)
	Find(key string) (string, error)
	Update(key string, value string)
	Save()
	DBIterator() *DBCusor
	Print()
}

func (dbc *DBCusor) Next() (k, v string) {
	if dbc.cursor >= len(dbc.keys) {
		return
	}
	k = dbc.keys[dbc.cursor].Interface().(string)

	v, _ = dbc.db.Find(k)

	dbc.cursor++

	return

}
