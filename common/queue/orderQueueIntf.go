package queue

type OrderQueueI interface {
	EnQueue(node interface{})
	DeQueue() interface{}
	Size() int
}
