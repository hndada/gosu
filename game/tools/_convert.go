package tools

import (
	"strconv"
)

func Atof(str string) float64 {
	float64Value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return float64Value
}

// always tries float64 conversion so far
func Atoi(str string) int {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return int(Atof(str))
	}
	return intValue
}

func Ftoa(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func Itoa(v int) string {
	return strconv.Itoa(v)
}
