package main

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/hndada/gosu"
)

func main() {
	var params struct {
		MusicDir string
		gosu.SearchParameter
	}
	_, err := toml.DecodeFile("q.toml", &params)
	if err != nil {
		panic(err)
	}
	existed := gosu.ChartSetList(params.MusicDir)
	banPath, err := filepath.Abs("ban.txt")
	if err != nil {
		panic(err)
	}
	ban := gosu.BanList(banPath)
	for _, r := range gosu.Search(params.SearchParameter) {
		if existed[r.SetId] {
			fmt.Printf("existed: %s\n", r.Filename())
		} else if ban[r.SetId] {
			fmt.Printf("banned: %s\n", r.Filename())
		} else {
			if r.Download(params.MusicDir) != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("downloaded: %s\n", r.Filename())
		}
	}

}
