package data

import (
	"fmt"
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

// Todo: should deleted charts be checked after Unmarshal ChartInfos?
// Todo: MessagePack when tags=release, JSON when tags=debug

func LoadData(fpath string, dst []any) {
	b, err := os.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		_ = os.Rename(fpath, fpath+".crashed") // Try rename if exists.
	}
	msgpack.Unmarshal(b, &dst)
}

// Todo: error notification?
func SaveData(fpath string, src []any) {
	b, err := msgpack.Marshal(&src)
	if err != nil {
		fmt.Printf("Failed to save data by an error: %s", err)
		return
	}
	err = os.WriteFile(fpath, b, 0644)
	if err != nil {
		fmt.Printf("Failed to save data by an error: %s", err)
		return
	}
}

// func LoadCharts(musicPath string) {
// 	const fname = "chart.db"
// 	b, err := os.ReadFile(fname)
// 	if err != nil {
// 		fmt.Println(err)
// 		// Rename will fail if not existed, but no have to handle error.
// 		os.Rename(fname, fname+".crashed")
// 	}
// 	msgpack.Unmarshal(b, &ChartInfos)
// 	LoadNewCharts(musicPath)
// }

// func SaveChartInfos() {
// 	b, err := msgpack.Marshal(&ChartInfos)
// 	if err != nil {
// 		fmt.Printf("Failed to save by an error: %s", err)
// 		return
// 	}
// 	err = os.WriteFile("chart.db", b, 0644)
// 	if err != nil {
// 		fmt.Printf("Failed to save by an error: %s", err)
// 		return
// 	}
// }
