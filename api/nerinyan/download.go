package nerinyan

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadBeatmapSet(id int, name string) {
	url := fmt.Sprintf("https://api.nerinyan.moe/d/%d", id)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/x-osu-beatmap-archive")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if name == "" {
		name = fmt.Sprintf("%d", id)
	}
	f, _ := os.Create(fmt.Sprintf("%s.osz", name))
	defer f.Close()
	f.Write(body)
}
