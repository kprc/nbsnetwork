package db

import (
	"bufio"
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"reflect"
	"fmt"
	"sync"
)

type filedbkev struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

type filedb struct {
	filepath string
	flock *sync.Mutex
	f        *os.File
	mkey     map[string]string
}

func NewFileDb(filepath string) NbsDbInter {
	return &filedb{filepath: filepath, mkey: make(map[string]string),flock:&sync.Mutex{}}
}

func (fdb *filedb)Print()  {
	for k,v:=range fdb.mkey{
		fmt.Println(k,v)
	}
}

func (fdb *filedb)DBIterator() *DBCusor {
	return &DBCusor{keys:reflect.ValueOf(fdb.mkey).MapKeys(),db:fdb}
}

func (fdb *filedb) Load() NbsDbInter {
	if fdb.filepath == "" {
		log.Fatal("No File ")
	}

	fdb.flock.Lock()
	defer fdb.flock.Unlock()

	flag := os.O_RDWR | os.O_APPEND

	if !tools.FileExists(fdb.filepath) {
		flag |= os.O_CREATE
	}


	if f, err := os.OpenFile(fdb.filepath, flag, 0755); err != nil {
		log.Fatal("Can't open file")
	} else {
		fdb.f = f
	}
	fdb.load()

	fdb.f.Close()

	fdb.f = nil

	return fdb
}

func (fdb *filedb) load() {
	if fdb.f == nil {
		return
	}
	bf := bufio.NewReader(fdb.f)

	for {
		if line, _, err := bf.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}

			if err == bufio.ErrBufferFull {
				log.Fatal("Buffer full")
				break
			}

			if len(line) > 0 {
				//pending drop it
				log.Fatal("Reading pending")
				break
			}

		} else {
			if len(line) > 0 {
				fdb.tomap(line)
			}
		}
	}
}

func (fdb *filedb) tomap(line []byte) {

	k := &filedbkev{}

	var delflag bool

	if line[0]=='-'{
		delflag = true
		line = line[1:]
	}

	if err := json.Unmarshal(line, k); err != nil {
		return
	} else {
		if delflag{
			delete(fdb.mkey,k.Key)
		}else{
			fdb.mkey[k.Key]= k.Value
		}
	}
}

func (fdb *filedb) Insert(key string, value string) error {
	if _, ok := fdb.mkey[key]; !ok {
		fdb.mkey[key] = value
		fdb.AppendSave(key,value,false)
	} else {
		return errors.New("Duplicate key")
	}

	return nil
}

func (fdb *filedb) Delete(key string) {
	v:=fdb.mkey[key]

	delete(fdb.mkey, key)
	fdb.AppendSave(key,v,true)
}

func (fdb *filedb) Find(key string) (string, error) {
	if v, ok := fdb.mkey[key]; ok {
		return v, nil
	}
	return "", errors.New("Not Found")
}

func (fdb *filedb) Update(key string, value string) {

	fdb.mkey[key] = value
	fdb.AppendSave(key,value,false)
}

func (fdb *filedb) write(data []byte) {
	if fdb.f == nil  {
		flag := os.O_WRONLY | os.O_TRUNC

		if !tools.FileExists(fdb.filepath) {
			flag |= os.O_CREATE
		}
		if f, err := os.OpenFile(fdb.filepath, flag, 0755); err != nil {
			log.Fatal("Can't open file")
			return
		} else {
			fdb.f = f
		}
	}

	fdb.f.Write(data)
}

func (fdb *filedb) Save() {

	fdb.flock.Lock()
	defer fdb.flock.Unlock()
	listkey := reflect.ValueOf(fdb.mkey).MapKeys()
	for _, key := range listkey {
		k := key.Interface().(string)

		fk := &filedbkev{}

		fk.Key = k
		fk.Value = fdb.mkey[k]

		if bj, err := json.Marshal(fk); err != nil {
			log.Println("save error", k, fk.Value)
		} else {
			fdb.write(bj)
			fdb.write([]byte("\r\n"))
		}

	}

	if fdb.f != nil {
		fdb.f.Close()
	}

	fdb.f = nil


}


func (fdb *filedb)AppendSave(key,value string,del bool)  {
	fdb.flock.Lock()
	defer fdb.flock.Unlock()

	if fdb.f == nil  {
		flag := os.O_WRONLY | os.O_APPEND

		if !tools.FileExists(fdb.filepath) {
			flag |= os.O_CREATE
		}
		if f, err := os.OpenFile(fdb.filepath, flag, 0755); err != nil {
			log.Fatal("Can't open file")
			return
		} else {
			fdb.f = f
		}
	}

	fk := &filedbkev{}

	fk.Key = key
	fk.Value = value

	if bj, err := json.Marshal(fk); err != nil {
		log.Println("save error", fk.Key, fk.Value)
	} else {
		if del {
			fdb.f.Write([]byte{'-'})
		}
		fdb.f.Write(bj)
		fdb.f.Write([]byte("\r\n"))
	}

	fdb.f.Close()
	fdb.f = nil

}
