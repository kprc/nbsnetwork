package hdb

import (
	"sync"
	"github.com/kprc/nbsnetwork/tools"
	"os"
	"log"
	"path"
	"bufio"
	"io"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"encoding/json"
	"errors"
	"reflect"
)

type HDBV struct {
	V string `json:"v"`
	Time int64 `json:"t"`
	Cnt int `json:"c"`
}

type FileHDBV struct {
	SaveLock sync.Mutex
	TimeFileName string
	Dbv FixedQueueIntf
	TotalCnt int
	hfdb *HistoryFileDB
}


type HistoryFileDB struct {
	MemHistoryCnt int
	SavePath string
	Mem map[string]*FileHDBV
	indexLock sync.Mutex
	indexFile *os.File
	indexFileName string
}

func New(cnt int,dbpath string) HistoryDBIntf  {
	if !tools.FileExists(dbpath){
		if err:=os.MkdirAll(dbpath,0755);err!=nil{
			log.Fatal(err)
		}
	}

	if cnt <=0 {
		cnt = 100
	}

	hfdb:=&HistoryFileDB{}
	hfdb.Mem = make(map[string]*FileHDBV)
	hfdb.SavePath = dbpath
	hfdb.MemHistoryCnt = cnt
	hfdb.indexFileName = path.Join(dbpath,"index.db")


	return hfdb
}

func (fv *FileHDBV)appendSave(v *HDBV)  {
	fv.SaveLock.Lock()
	defer fv.SaveLock.Unlock()

	flag:=os.O_WRONLY|os.O_APPEND
	if !tools.FileExists(fv.TimeFileName){
		flag |= os.O_CREATE
	}
	if f,err:=os.OpenFile(fv.TimeFileName,flag,0755);err!=nil{
		log.Fatal("can't open file",fv.TimeFileName)
		return
	}else{
		if bj,err:=json.Marshal(*v);err!=nil{
			log.Println("save error",v.V)
		}else{
			f.Write(bj)
			f.Write([]byte("\r\n"))
		}
		f.Close()
	}

}

func (hfdb *HistoryFileDB)loadIndex()  {
	if hfdb.indexFile == nil{
		return
	}

	bf:=bufio.NewReader(hfdb.indexFile)

	for{
		if line,_,err:=bf.ReadLine();err!=nil{
			if err == io.EOF{
				break
			}
			if err == bufio.ErrBufferFull{
				log.Fatal("Buffer full")
			}
			if len(line) > 0{
				log.Fatal("Reading pending")
				break
			}
		}else{
			if len(line) > 0{
				hfdb.Mem[string(line)] = &FileHDBV{TimeFileName:hfdb.GetTimeFileName(string(line)),hfdb:hfdb}
			}
		}
	}
}

func (hfdb *HistoryFileDB)GetTimeFileName(key string) string  {
	hash:=sha256.Sum256([]byte(key))
	hstr := base58.Encode(hash[:])

	lhstr := len(hstr)
	return path.Join(hfdb.SavePath,"historyFile",
		hstr[lhstr-8:lhstr-6],hstr[lhstr-6:lhstr-4],hstr[lhstr-4:lhstr-2],hstr[lhstr-2:], hstr)
}

func (fv *FileHDBV)loadFile()  {
	fv.SaveLock.Lock()
	defer fv.SaveLock.Unlock()

	flag := os.O_RDWR | os.O_APPEND

	if !tools.FileExists(fv.TimeFileName) {
		flag |= os.O_CREATE
	}

	var (
		f *os.File
		err error
		)
	

	f, err = os.OpenFile(fv.TimeFileName, flag, 0755) 
	if err != nil {
		log.Fatal("Can't open file")
	} 
	defer f.Close()	

	fv.Dbv = NewFixedQueue(fv.hfdb.MemHistoryCnt, func(v1 interface{}, v2 interface{}) int {
		d1,d2:=v1.(*HDBV),v2.(*HDBV)
		if d1.V == d2.V && d1.Time == d2.Time{
			return 0
		}else{
			return 1
		}
	})

	bf:=bufio.NewReader(f)

	for{
		if line,_,err:=bf.ReadLine();err!=nil{
			if err == io.EOF{
				break
			}

			if err == bufio.ErrBufferFull{
				log.Fatal("Buffer full")
				break
			}

			if len(line) > 0{
				log.Fatal("Reading pending")
				break
			}

		}else{
			if len(line) > 0{
				dbv:=&HDBV{}
				err:=json.Unmarshal([]byte(line),dbv)
				if err!=nil{
					log.Println("unmarshall message failed",line)
					continue
				}
				fv.Dbv.EnQueue(dbv)
				fv.TotalCnt ++
			}
		}
	}
}

func (hfdb *HistoryFileDB)loadValue()  {
	for k,_:=range hfdb.Mem{
		v := hfdb.Mem[k]
		v.loadFile()
	}
}


func (hfdb *HistoryFileDB)Load() HistoryDBIntf {
	hfdb.indexLock.Lock()
	defer hfdb.indexLock.Unlock()

	flag:=os.O_RDWR | os.O_APPEND

	if !tools.FileExists(hfdb.indexFileName){
		flag |= os.O_CREATE
	}

	if f,err:=os.OpenFile(hfdb.indexFileName,flag,0755);err!=nil{
		log.Fatal("Can't open file "+hfdb.indexFileName)
	}else{
		hfdb.indexFile = f
	}

	hfdb.loadIndex()
	hfdb.loadValue()

	hfdb.indexFile.Close()
	hfdb.indexFile = nil

	return hfdb
}

func (hfdb *HistoryFileDB)insertIndex(key string)  {
	if _,ok:=hfdb.Mem[key];ok{
		return
	}

	hfdb.Mem[key] = &FileHDBV{TimeFileName:hfdb.GetTimeFileName(key),hfdb:hfdb}
	hfdb.appendInex(key,false)
}

func (hfdb *HistoryFileDB)appendInex(key string,del bool)  {
	hfdb.indexLock.Lock()
	defer hfdb.indexLock.Unlock()

	if hfdb.indexFile == nil{
		flag := os.O_WRONLY | os.O_APPEND

		if !tools.FileExists(hfdb.indexFileName) {
			flag |= os.O_CREATE
		}
		if f, err := os.OpenFile(hfdb.indexFileName, flag, 0755); err != nil {
			log.Fatal("Can't open file")
			return
		} else {
			hfdb.indexFile = f
		}
	}

	if del{
		hfdb.indexFile.Write([]byte{'-'})
	}

	hfdb.indexFile.Write([]byte(key))
	hfdb.indexFile.Write([]byte("\r\n"))

	hfdb.indexFile.Close()
	hfdb.indexFile = nil

}


func (hfdb *HistoryFileDB)Insert(key string,value string) error  {
	_,ok:=hfdb.Mem[key]
	if !ok{
		hfdb.insertIndex(key)
	}

	dbv:=hfdb.Mem[key]
	v := &HDBV{V:value,Time:tools.GetNowMsTime(),Cnt:dbv.TotalCnt}
	dbv.appendSave(v)
	if dbv.Dbv == nil{
		dbv.Dbv = NewFixedQueue(hfdb.MemHistoryCnt, func(v1 interface{}, v2 interface{}) int {
			d1,d2:=v1.(*HDBV),v2.(*HDBV)
			if d1.V == d2.V && d1.Time == d2.Time{
				return 0
			}else{
				return 1
			}
		})
	}
	dbv.Dbv.EnQueue(v)
	dbv.TotalCnt ++

	return nil
}

func (hfdb *HistoryFileDB)Delete(key string) {
	if v,ok:=hfdb.Mem[key];!ok{
		return
	}else{
		v.SaveLock.Lock()
		defer v.SaveLock.Unlock()
		os.Remove(v.TimeFileName)
	}

	hfdb.appendInex(key,true)

	delete(hfdb.Mem,key)

	return

}

func (hfdb *HistoryFileDB) write(data []byte) {
	if hfdb.indexFile == nil  {
		flag := os.O_WRONLY | os.O_TRUNC

		if !tools.FileExists(hfdb.indexFileName) {
			flag |= os.O_CREATE
		}
		if f, err := os.OpenFile(hfdb.indexFileName, flag, 0755); err != nil {
			log.Fatal("Can't open file")
			return
		} else {
			hfdb.indexFile = f
		}
	}

	hfdb.indexFile.Write(data)
}



func (hfdb *HistoryFileDB)Save()  {
	hfdb.indexLock.Lock()
	defer hfdb.indexLock.Unlock()

	for k,_:=range hfdb.Mem{
		hfdb.write([]byte(k))
		hfdb.write([]byte("\r\n"))
	}

	hfdb.indexFile.Close()
	hfdb.indexFile = nil
}

func (hfdb *HistoryFileDB)Find(key string,start int, topn int) ([]*HDBV,error)  {

	if start < 0{
		start = 0
	}

	if topn <= 0{
		topn = hfdb.MemHistoryCnt
	}

	v,ok:=hfdb.Mem[key]
	if !ok{
		return nil,errors.New("no key in db")
	}

	if start > v.TotalCnt{
		return nil,errors.New("start pos error")
	}

	spos:=v.TotalCnt - hfdb.MemHistoryCnt
	if spos > 0{
		if start < spos {
			start = spos
		}
	}

	arr:=v.Dbv.GetTopN(start,topn)

	var as []*HDBV

	for _,i:=range arr{
		as = append(as,i.(*HDBV))
	}

	return as,nil
}

func (hfdb *HistoryFileDB)FindBlock(key string) (*FileHDBV,error)  {
	v,ok:=hfdb.Mem[key]
	if !ok{
		return nil,errors.New("not found")
	}
	return v,nil
}


func (hfdb *HistoryFileDB)DBIterator() *DBCusor  {
	return &DBCusor{keys:reflect.ValueOf(hfdb.Mem).MapKeys(),hdb:hfdb}
}
