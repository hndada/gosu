package beatmap

import (
	"strconv"
	"strings"
)

type Info map[string]interface{}

func (m Info) PutStr(kv []string) { m[kv[0]] = kv[1] }

func (m Info) PutInt(kv []string) error {
	d, err := strconv.Atoi(kv[1])
	if err != nil {
		return err
	}
	m[kv[0]] = d
	return nil
}

func (m Info) PutF64(kv []string) error {
	f, err := strconv.ParseFloat(kv[1], 64)
	if err != nil {
		return err
	}
	m[kv[0]] = f
	return nil
}

func (m Info) PutBool(kv []string) error {
	b, err := strconv.Atoi(kv[1])
	if err != nil {
		return err
	}
	switch b {
	case 1:
		m[kv[0]] = true
	default: // todo: check how osu! handles the invalid value
		m[kv[0]] = false
	}
	return nil
}

func (m Info) PutIntSlice(kv []string) error {
	sSlice := strings.Split(kv[1], ",")
	dSlice := make([]int, 0, len(sSlice))
	for _, s := range sSlice {
		d, err := strconv.Atoi(s)
		if err != nil {
			return err
		} else {
			dSlice = append(dSlice, d)
		}
	}
	m[kv[0]] = dSlice
	return nil
}

func (m Info) PutStrSlice(kv []string) { m[kv[0]] = strings.Split(kv[1], ",") }
