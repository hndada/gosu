package choose

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ModePiano = 3 // osu!mania
	ModeDrum  = 1 // osu!taiko
)

const (
	StatusGraveyard = -2
	StatusWIP       = -1
	StatusPending   = iota
	StatusRanked
	StatusApproved
	StatusQualified
	StatusLoved
)
const (
	Unranked = iota
	Ranked
	Approved
	Qualified
	Loved
)

type SearchParam struct {
	Query   string
	Mode    int
	SubMode int

	page int
}
type ChartSet struct {
	SetId            int
	ChildrenBeatmaps []Chart
	RankedStatus     int
	ApprovedDate     string
	LastUpdate       string
	LastChecked      string
	Artist           string
	Title            string
	Creator          string
	Source           string
	Tags             string
	HasVideo         bool
	Genre            int
	Language         int
	Favourites       int
	Disabled         int
}
type Chart struct {
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

const amount = 25

func (p SearchParam) URL() *url.URL {
	u, err := url.Parse("https://api.chimu.moe/search")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	vs := url.Values{}
	vs.Add("query", p.Query)
	vs.Add("mode", strconv.Itoa(p.Mode))
	vs.Add("min_cs", strconv.Itoa(p.SubMode))
	vs.Add("max_cs", strconv.Itoa(p.SubMode))
	vs.Add("amount", strconv.Itoa(amount))
	vs.Add("offset", strconv.Itoa(p.page*amount))
	u.RawQuery = vs.Encode()
	return u
}
func (p *SearchParam) Search() (sets []ChartSet, err error) {
	u := p.URL()
	fmt.Printf("Search page %d\n", p.page)
	resp, err := http.Get(u.String())
	if err != nil || resp.StatusCode == 404 {
		return sets, err
	}
	j, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return sets, err
	}
	result := struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Data    []ChartSet `json:"data"`
	}{
		Data: make([]ChartSet, 0, amount),
	}
	err = json.Unmarshal(j, &result)
	if err != nil {
		return sets, err
	}
	if len(result.Data) == 0 {
		return sets, err
	}
	sets = append(sets, result.Data...)
	p.page++
	return
}

func (c Chart) Select() (fsys fs.FS, name string, err error) {
	// const noVideo = 1
	// u := fmt.Sprintf("%s%d?n=%d", APIDownload, c.ParentSetId, noVideo)
	u := c.URLDownload()
	fmt.Printf("download URL: %s\n", u)
	// err will be assigned to return value 'err'.
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fsys, err = zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return
	}
	return fsys, c.OsuFile, err
}

// https://osu.ppy.sh/docs/index.html#beatmapsetcompact-covers
// slimcover: 1920x360
// cover: 900x250
// card: 400x140
// list: 150x150
const APIBeatmap = "https://assets.ppy.sh/beatmaps"
const Large = "@2x"

func (c ChartSet) URLCover(kind, suffix string) string {
	return fmt.Sprintf("%s/%d/covers/%s%s.jpg", APIBeatmap, c.SetId, kind, suffix)
}
func (c ChartSet) URLPreview() string {
	return fmt.Sprintf("b.ppy.sh/preview/%d.mp3", c.SetId)
}
func (c ChartSet) URLDownload() string {
	return fmt.Sprintf("https://api.chimu.moe/v1/d/%d", c.SetId)
}
func (c Chart) URLDownload() string {
	return fmt.Sprintf("https://api.chimu.moe/v1/%s", c.DownloadPath)
}
