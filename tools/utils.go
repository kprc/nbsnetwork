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