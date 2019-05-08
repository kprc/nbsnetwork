package tools

import "time"

func ResizeHash(hash uint) uint {
	cnt := 0
	for{
		hash = hash >>1
		if hash > 0{
			cnt ++
		}else {
			break
		}
	}
	return 1<<uint(cnt)
}


func GetNowMsTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetRealPos(pos uint64) uint64  {
	return (pos & 0xFFFFFFFF)
}

func AssemblePos(pos uint64,typ uint32) uint64 {
	var typ1 uint64

	typ1 = uint64(typ)

	typ1 = typ1 << 32

	pos = pos & 0xFFFFFFFF

	return typ1 | pos

}

func GetTypFromPos(pos uint64) uint32  {
	typ := (pos >> 32) & 0xFFFFFFFF

	return uint32(typ)
}

