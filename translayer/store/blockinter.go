package store


type BlockInter interface {
	GetSn() uint64
}



var FBlockHash = func(v interface{}) uint {
	blk:=v.(BlockInter)

	return uint(blk.GetSn()&0x7F)
}

var FBlockEquals = func(v1 interface{},v2 interface{}) int{
	blk1:=v1.(BlockInter)
	blk2:=v2.(BlockInter)

	if blk1.GetSn() == blk2.GetSn() {
		return 0
	}

	return 1
}

