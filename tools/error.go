package tools

import (
	"errors"
	"os"
	"strings"
)

type ValError struct {
	Name string
	Val  string
	Err  error
}

type ValsError struct {
	Name string
	Vals []string
	Err  error
}

var ErrSyntax = errors.New("invalid value")
var ErrFlow = errors.New("unintended logic flow")

func (e *ValError) Error() string {
	return e.Name + ": " + e.Val + ": " + e.Err.Error()
}

func (e *ValsError) Error() string {
	return e.Name + ": " + strings.Join(e.Vals, ", ") + ": " + e.Err.Error()
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
func CheckDel(err error, fPath string) {
	if err != nil {
		err = os.Remove(fPath)
		Check(err)
	}
}
