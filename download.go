package gosu

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SearchParameter struct {
	Query       string
	Status      int
	Mode        int
	MinKeyCount int
	MaxKeyCount int
	MinOsuSR    float64
	MaxOsuSR    float64
	MinLength   int
	MaxLength   int
}

const (
	unranked = iota
	ranked
	approved
	qualified
	loved
)
const amount = 5

var offset int
var (
	exist map[int]bool
	ban   map[int]bool
)

const (
	chimuURL         = "https://api.chimu.moe/v1/"
	chimuURLSearch   = chimuURL + "search"
	chimuURLDownload = chimuURL + "download/"
)

type ChimuResult struct {
	SetId            int
	ChildrenBeatmaps []struct {
		BeatmapId        int
		ParentSetId      int
		DiffName         string
		FileMD5          string
		Mode             int
		BPM              float64
		AR               float64
		OD               float64
		CS               float64
		HP               float64
		TotalLength      int
		HitLength        int
		Playcount        int
		Passcount        int
		MaxCombo         int
		DifficultyRating float64
		OsuFile          string
		DownloadPath     string
	}
	RankedStatus int
	ApprovedDate string
	LastUpdate   string
	LastChecked  string
	Artist       string
	Title        string
	Creator      string
	Source       string
	Tags         string
	HasVideo     bool
	Genre        int
	Language     int
	Favourites   int
	Disabled     bool
}

// set searchParameter value
func Search(params SearchParameter) {
	u, err := url.Parse(chimuURLSearch)
	if err != nil {
		panic(err)
	}
	vs := url.Values{}
	vs.Add("query", params.Query)
	vs.Add("amount", strconv.Itoa(amount))
	vs.Add("offset", "0")
	vs.Add("status", strconv.Itoa(params.Status))
	vs.Add("mode", strconv.Itoa(params.Mode))
	vs.Add("min_cs", strconv.Itoa(params.MinKeyCount))
	vs.Add("max_cs", strconv.Itoa(params.MaxKeyCount))
	vs.Add("min_diff", strconv.FormatFloat(params.MinOsuSR, 'f', -1, 64))
	vs.Add("max_diff", strconv.FormatFloat(params.MaxOsuSR, 'f', -1, 64))
	vs.Add("min_length", strconv.Itoa(params.MinLength))
	vs.Add("max_length", strconv.Itoa(params.MaxLength))
	u.RawQuery = vs.Encode()

	for {
		resp, err := http.Get(u.String())
		if err != nil {
			panic(err)
		} else if resp.StatusCode == 404 {
			fmt.Printf("URL: %s\n", u.String())
			fmt.Println("no more result")
			os.Exit(1)
		}
		j, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}
		result := struct {
			Code    int           `json:"code"`
			Message string        `json:"message"`
			Data    []ChimuResult `json:"data"`
		}{}
		result.Data = make([]ChimuResult, 0, amount)
		err = json.Unmarshal(j, &result)
		if err != nil {
			panic(err)
		}
		// if len(result.Data) == 0 {
		// 	fmt.Println("no data")
		// 	break
		// }
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Printf("cwd: %s\n", cwd)
		for _, r := range result.Data {
			if exist[r.SetId] {
				fmt.Printf("existed: %s\n", r.filename())
			} else if ban[r.SetId] {
				fmt.Printf("banned: %s\n", r.filename())
			} else {
				if r.download() != nil {
					panic(err)
				}
				fmt.Printf("downloaded: %s\n", r.filename())
			}
		}
		offset += amount
		vs.Set("offset", strconv.Itoa(offset))
		u.RawQuery = vs.Encode()
	}
}

func (r ChimuResult) download() error {
	const noVideo = 1
	u := fmt.Sprintf("%s%d?n=%d", chimuURLDownload, r.SetId, noVideo)
	fmt.Printf("download URL: %s\n", u)
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filepath.Join(cwd, r.filename()))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
func (r ChimuResult) filename() string {
	name := fmt.Sprintf("%d %s - %s.osz", r.SetId, r.Artist, r.Title)
	for _, letter := range []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"} {
		name = strings.ReplaceAll(name, letter, "-")
	}
	return name
}
