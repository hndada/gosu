package tools

// NoFound is the const for an element that is not found yet
// but its default value should not be zero.
const NoFound = -1
const (
	left  = -1
	alter = 0
	right = 1
)

// GetIntSlice function returns int slice,
// filled with given default value.
func GetIntSlice(len int, defaultValue int) []int {
	newSlice := make([]int, len)
	for i := range newSlice {
		newSlice[i] = defaultValue
	}
	return newSlice
}

func Neighbors(slice []int, i int) [2]int {
	ns := [2]int{NoFound, NoFound}
	uBound := len(slice)

	var cursor, v int
	for ni, direct := range [2]int{left, right} {
		for offset := 1; ; offset++ {
			cursor = i + offset*direct
			if cursor < 0 || cursor >= uBound {
				break
			}
			v = slice[cursor]
			if v == NoFound {
				continue
			}
			ns[ni] = v
			break
		}
	}
	return ns
}
