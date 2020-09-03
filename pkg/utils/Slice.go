package utils

func IntIndexOf(IntSlice *[]int, Num int) int {
	for i, member := range *IntSlice {
		if member == Num {
			return i
		}
	}
	return -1
}

func ChoiceSlice(IntSlice *[]int) int {
	_IntSlice := *IntSlice
	if len(_IntSlice) == 0 {
		return 0
	} else if len(_IntSlice) == 1 {
		return _IntSlice[0]
	}
	// 返回0~n-1之间的随机数
	return _IntSlice[RandInt(0, len(*IntSlice)-1)]
}
