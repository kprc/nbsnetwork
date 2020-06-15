package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/hdb"
	"github.com/kprc/nbsnetwork/tools"
	"path"
)

func main() {

	home, _ := tools.Home()
	testdbpath := path.Join(home, "testhdb")

	db := hdb.New(5, testdbpath).Load()

	//db.Insert("a","11111")
	//db.Insert("a","22222")
	//db.Insert("a","33333")
	//db.Insert("a","44444")
	//db.Insert("a","55555")
	//db.Insert("a","66666")
	//db.Insert("a","77777")
	//db.Insert("a","88888")
	//db.Insert("a","99999")
	//
	//db.Insert("b","11111")
	//db.Insert("b","22222")
	//db.Insert("b","33333")
	//db.Insert("b","44444")
	//db.Insert("b","55555")
	//db.Insert("b","66666")
	//db.Insert("b","77777")
	//db.Insert("b","88888")
	//db.Insert("b","99999")
	//
	//db.Insert("c","11111")
	//db.Insert("c","22222")
	//db.Insert("c","33333")
	//db.Insert("c","44444")
	//db.Insert("c","55555")
	//db.Insert("c","66666")
	//db.Insert("c","77777")
	//db.Insert("c","88888")
	//db.Insert("c","99999")

	//db.Delete("a")

	dbc := db.DBIterator()
	for {
		k, v := dbc.Next()

		if k == "" {
			break
		}

		fmt.Println(k)
		for i := 0; i < len(v); i++ {
			vv := v[i]

			fmt.Println(vv.Cnt, vv.V)
		}

	}

	hs, err := db.FindMem("a", 0, 20)
	if err != nil {
		fmt.Println("err: ", err)
	}

	for i := 0; i < len(hs); i++ {
		vv := hs[i]

		fmt.Println(vv.Cnt, vv.V)
	}

}
