package main

import (
	"github.com/BurntSushi/toml"
	"github.com/hndada/gosu"
)

func main() {
	var params gosu.SearchParameter
	_, err := toml.DecodeFile("q.toml", &params)
	if err != nil {
		panic(err)
	}
	gosu.Search(params)
}
