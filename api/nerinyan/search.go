package nerinyan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// func LocalParse() {
// 	// beatmapSets := make([]BeatmapSet, 100)
// 	json.Unmarshal([]byte(result), &beatmapSets)

// 	jsonBytes, err := json.Marshal(beatmapSets)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	f, _ := os.Create("beatmapsets.json")
// 	defer f.Close()
// 	f.Write(jsonBytes)
// }

// Get all beatmap info
func Search() {
	var beatmapSetsAll []BeatmapSet
	baseURL := "https://api.nerinyan.moe/search"
	page := 0
	var pageSize int = 1e3
	var statusString string
	for s := -2; s <= 4; s++ {
		statusString += fmt.Sprintf("%d,", s)
	}
	statusString = statusString[:len(statusString)-1]

	// make dir with the current time
	dir := time.Now().Format("20060102150405")
	os.Mkdir(dir, 0777)
	for {
		url := fmt.Sprintf("%s?s=%s&m=1,3&p=%d&ps=%d", baseURL, statusString, page, pageSize)
		fmt.Printf("Fetching %s\n", url)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		f, _ := os.Create(fmt.Sprintf("%s/%d.json", dir, page))
		defer f.Close()
		f.Write(body)

		var beatmapSets []BeatmapSet
		json.Unmarshal(body, &beatmapSets)
		if len(beatmapSets) == 0 {
			break
		}

		beatmapSetsAll = append(beatmapSetsAll, beatmapSets...)
		page++
		time.Sleep(600 * time.Millisecond)
	}

	// output to json
	jsonBytes, _ := json.Marshal(beatmapSetsAll)
	f, _ := os.Create("all.json")
	defer f.Close()
	f.Write(jsonBytes)

}
