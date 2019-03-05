package tools

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
