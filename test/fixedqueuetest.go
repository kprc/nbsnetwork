package main

import (
	"fmt"
	"github.com/kprc/nbsnetwork/hdb"
)

func main() {
	fq := hdb.NewFixedQueue(3, func(v1 interface{}, v2 interface{}) int {
		d1, d2 := v1.(int), v2.(int)
		if d1 == d2 {
			return 0
		} else {
			return 1
		}

	})

	fq.EnQueue(1)
	fq.EnQueue(2)
	fq.EnQueue(3)
	fq.EnQueue(4)
	fq.EnQueue(5)

	cusor := fq.Iterator()
	for {
		x := cusor.Next()
		if x == nil {
			break
		}

		fmt.Println(x.(int))
	}

}
