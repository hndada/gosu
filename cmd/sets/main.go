package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/hndada/gosu/api/nerinyan"
)

func main() {
	// nerinyan.FetchBeatmapInfos()
	// read all.json then unmarshal
	//
	f, _ := os.Open("all.json")
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var all []nerinyan.BeatmapSet
	json.Unmarshal(data, &all)
	// var (
	// 	pianoRanked     []nerinyan.BeatmapSet
	// 	pianoSemiRanked []nerinyan.BeatmapSet
	// 	pianoUnranked   []nerinyan.BeatmapSet
	// 	drumRanked      []nerinyan.BeatmapSet
	// 	drumSemiRanked  []nerinyan.BeatmapSet
	// 	drumUnranked    []nerinyan.BeatmapSet
	// )
	var (
		pianoRanked     []int
		pianoSemiRanked []int
		pianoUnranked   []int
		drumRanked      []int
		drumSemiRanked  []int
		drumUnranked    []int
	)
	for _, s := range all {
		var isPiano, isDrum bool
		for _, d := range s.Beatmaps {
			switch d.ModeInt {
			case 1:
				isPiano = true
			case 3:
				isDrum = true
			}
		}
		if isPiano {
			switch s.Ranked {
			case 1, 3:
				pianoRanked = append(pianoRanked, s.ID)
			case 2, 4:
				pianoSemiRanked = append(pianoSemiRanked, s.ID)
			default:
				pianoUnranked = append(pianoUnranked, s.ID)
			}
		}
		if isDrum {
			switch s.Ranked {
			case 1, 3:
				drumRanked = append(drumRanked, s.ID)
			case 2, 4:
				drumSemiRanked = append(drumSemiRanked, s.ID)
			default:
				drumUnranked = append(drumUnranked, s.ID)
			}
		}
	}

	slices := map[string][]int{
		"pianoRanked":     pianoRanked,
		"pianoSemiRanked": pianoSemiRanked,
		"pianoUnranked":   pianoUnranked,
		"drumRanked":      drumRanked,
		"drumSemiRanked":  drumSemiRanked,
		"drumUnranked":    drumUnranked,
	}
	const dir = "list"
	os.Mkdir(dir, 0777)
	for name, slice := range slices {
		// Convert the slice to a string representation
		data := sliceToString(slice)

		// Write the data to a file
		fileName := dir + "/" + name + ".txt"
		err := ioutil.WriteFile(fileName, []byte(data), 0644)
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Slice data saved to %s\n", fileName)
		}
	}
}
func sliceToString(slice []int) string {
	var strSlice []string
	for _, value := range slice {
		strSlice = append(strSlice, strconv.Itoa(value))
	}
	return strings.Join(strSlice, ", ")
}
