package osu

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func keyValue(line, delimiter string) (key, value string, err error) {
	kv := strings.SplitN(line, delimiter, 2)
	if len(kv) < 2 {
		return "", "", fmt.Errorf("%s: key value not enough length", line)
	}

	k := strings.TrimSpace(kv[0])
	// TrimRightFunc is preferred to TrimSpace, because
	// TrimSpace may also trim some values starting with space.
	v := strings.TrimRightFunc(kv[1], unicode.IsSpace)
	return k, v, nil
}

func parseInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return int(f), nil
	}
	return i, nil
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func parsePoint(s string) (p [2]int, err error) {
	xy := strings.Split(s, `:`)
	if len(xy) < 2 {
		return p, fmt.Errorf("point has not enough length: %s", s)
	}
	for i := 0; i < 2; i++ {
		if p[i], err = parseInt(xy[i]); err != nil {
			return
		}
	}
	return
}
