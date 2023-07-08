package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/hndada/gosu/api/nerinyan"
)

func main() {
	filenames := []string{"pianoRanked", "pianoSemiRanked", "drumRanked", "drumSemiRanked"}
	for _, dir := range filenames {
		os.Mkdir(dir, 0777)
		f, _ := os.Open(fmt.Sprintf("%s.txt", dir))
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		ids := strings.Split(string(data), ", ")
		for _, idString := range ids {
			id, _ := strconv.Atoi(idString)
			fmt.Printf("Downloading %d\n", id)
			name := fmt.Sprintf("%s/%d", dir, id)
			nerinyan.DownloadBeatmapSet(id, name) // id goes name when name is empty
		}
	}
}
