package bus


type FSendData func(v interface{}, uid string)


type busblock struct {
	data interface{}
	uid string
	f FSendData
}

type BusBlock interface {
	SetData(data interface{})
	GetData() interface{}
	SetUid(uid string)
	GetUid() string
	SetFSendData(f FSendData)
	GetFSendData() FSendData
}

func NewBusBlock() BusBlock {
	return &busblock{}
}


func (bb *busblock)SetData(data interface{})  {
	bb.data = data
}

func (bb *busblock)GetData() interface{}  {
	return bb.data
}

func (bb *busblock)SetUid(uid string)  {
	bb.uid  = uid
}

func (bb *busblock)GetUid() string  {
	return bb.uid
}

func (bb *busblock)SetFSendData(f FSendData)  {
	bb.f = f
}

func (bb *busblock)GetFSendData() FSendData  {
	return bb.f
}

func SendBusBlock(v interface{}) {
	bb:=v.(BusBlock)
	if bb.GetFSendData() !=nil{
		bb.GetFSendData()(bb.GetData(),bb.GetUid())
	}
}

