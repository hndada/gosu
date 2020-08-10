package tools

import (
	"strconv"
	"strings"
)

func Atoi(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return int(f), nil
	}
	return v, nil
}

func PairInt(s, sep string) ([2]int, error) {
	var pair [2]int
	vs := strings.Split(s, sep)
	for i := range pair {
		v, err := Atoi(vs[i])
		if err != nil {
			return pair, err
		}
		pair[i] = v
	}
	return pair, nil
}