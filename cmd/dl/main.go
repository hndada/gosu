package main

import (
	"github.com/hndada/gosu"
)

func main() {
	params := gosu.SearchParameter{
		Query:       "miiro",
		Status:      1,
		Mode:        3,
		MinKeyCount: 4,
		MaxKeyCount: 4,
		MinOsuSR:    2.5,
		MaxOsuSR:    4,
		MinLength:   60,
		MaxLength:   120,
	}
	gosu.Search(params)
}
