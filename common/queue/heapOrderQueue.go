package queue

type HeapType int

const (
	TopMax HeapType = 0
	TopMin HeapType = 1
)

type HeapOrderQueue struct {
	data []interface{}
	cmp  func(v1, v2 interface{}) int
	typ  HeapType
}

func NewHeapOrderQueue(data []interface{}, cmp func(v1, v2 interface{}) int, typ HeapType) *HeapOrderQueue {
	return &HeapOrderQueue{
		data: data,
		cmp:  cmp,
		typ:  typ,
	}
}

func (hoq *HeapOrderQueue) Dup() *HeapOrderQueue {
	d := &HeapOrderQueue{
		cmp: hoq.cmp,
		typ: hoq.typ,
	}

	data := make([]interface{}, len(hoq.data))

	for i := 0; i < len(hoq.data); i++ {
		data[i] = hoq.data[i]
	}

	d.data = data

	return d
}

func (hoq *HeapOrderQueue) DeQueue() interface{} {

	l := hoq.Size()

	if l == 0 {
		return nil
	}
	d := hoq.data[0]
	if l == 1 {
		hoq.data = nil
		return d
	}

	hoq.data[0] = hoq.data[l-1]
	hoq.data = hoq.data[:l-1]

	hoq.shiftDown(0)

	return d
}

func (hoq *HeapOrderQueue) Peek() interface{} {

	l := hoq.Size()

	if l == 0 {
		return nil
	}

	return hoq.data[0]

}

func (hoq *HeapOrderQueue) EnQueue(v interface{}) {

	hoq.data = append(hoq.data, v)

	hoq.shiftUP(hoq.Size() - 1)

	return
}

func (hoq *HeapOrderQueue) Build() *HeapOrderQueue {

	data := hoq.data
	hoq.data = nil

	for i := 0; i < len(data); i++ {
		hoq.data = append(hoq.data, data[i])
		hoq.shiftUP(i)
	}
	return hoq
}

func (hoq *HeapOrderQueue) swap(i, j int) {
	t := hoq.data[i]
	hoq.data[i] = hoq.data[j]
	hoq.data[j] = t
}

func (hoq *HeapOrderQueue) shiftUP(i int) {
	idx := i
	for {
		if idx == 0 {
			break
		}
		pidx := (idx+1)/2 - 1

		if hoq.typ == TopMin {
			if hoq.cmp(hoq.data[pidx], hoq.data[idx]) <= 0 {
				break
			}
		}
		if hoq.typ == TopMax {
			if hoq.cmp(hoq.data[pidx], hoq.data[idx]) >= 0 {
				break
			}
		}

		hoq.swap(idx, pidx)

		idx = pidx

	}
}

func (hoq *HeapOrderQueue) shiftDown(i int) {
	pidx := i
	for {
		left := 2*(pidx+1) - 1
		right := 2 * (pidx + 1)

		cmpidx := left

		if hoq.typ == TopMin {
			if left >= hoq.Size() {
				break
			}

			if right < hoq.Size() {
				if hoq.cmp(hoq.data[left], hoq.data[right]) > 0 {
					cmpidx = right
				}
			}

			if hoq.cmp(hoq.data[pidx], hoq.data[cmpidx]) <= 0 {
				break
			}
		}
		if hoq.typ == TopMax {
			if left >= hoq.Size() {
				break
			}
			if right < hoq.Size() {
				if hoq.cmp(hoq.data[left], hoq.data[right]) < 0 {
					cmpidx = right
				}
			}

			if hoq.cmp(hoq.data[pidx], hoq.data[cmpidx]) >= 0 {
				break
			}
		}

		hoq.swap(pidx, cmpidx)

		pidx = cmpidx

	}
}

func (hoq *HeapOrderQueue) Size() int {
	return len(hoq.data)
}

func (hoq *HeapOrderQueue) orderArray() []interface{} {
	var data []interface{}

	for {
		v := hoq.DeQueue()
		if v == nil {
			return data
		}
		data = append(data, v)
	}
}

func (hoq *HeapOrderQueue) orderArray2() []interface{} {
	var data = make([]interface{}, hoq.Size())
	idx := hoq.Size() - 1
	for {
		v := hoq.DeQueue()
		if v == nil {
			return data
		}
		data[idx] = v
		idx--
	}
}

func (hoq *HeapOrderQueue) ArrayASC() []interface{} {
	if hoq.typ == TopMin {
		return hoq.orderArray()
	} else {
		return hoq.orderArray2()
	}
}

func (hoq *HeapOrderQueue) ArrayDESC() []interface{} {
	if hoq.typ == TopMax {
		return hoq.orderArray()
	} else {
		return hoq.orderArray2()
	}
}

func (hoq *HeapOrderQueue) ArrayASC2() []interface{} {

	dup := hoq.Dup()

	if dup.typ == TopMin {
		return dup.orderArray()
	} else {
		return dup.orderArray2()
	}
}

func (hoq *HeapOrderQueue) ArrayDESC2() []interface{} {

	dup := hoq.Dup()

	if dup.typ == TopMax {
		return dup.orderArray()
	} else {
		return dup.orderArray2()
	}
}
