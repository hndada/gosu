package tools

import "strconv"

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
