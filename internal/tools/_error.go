package tools

import (
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

func (e *ValError) Error() string {
	return e.Name + ": " + e.Val + ": " + e.Err.Error()
}

func (e *ValsError) Error() string {
	return e.Name + ": " + strings.Join(e.Vals, ", ") + ": " + e.Err.Error()
}
