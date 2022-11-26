package choose

import (
	"encoding/json"
	"io"
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
const amount = 25

func (s *Scene) LoadChartSetList() (err error) {
	u, err := url.Parse("https://api.chimu.moe/search")
	if err != nil {
		return
	}
	vs := url.Values{}
	vs.Add("query", s.Query.Text)
	vs.Add("mode", strconv.Itoa(s.mode))
	vs.Add("min_cs", strconv.Itoa(s.subMode))
	vs.Add("max_cs", strconv.Itoa(s.subMode))
	vs.Add("amount", strconv.Itoa(amount))
	vs.Add("offset", strconv.Itoa(s.page*amount))
	u.RawQuery = vs.Encode()
	resp, err := http.Get(u.String())
	if err != nil || resp.StatusCode == 404 {
		return
	}
	defer resp.Body.Close()
	j, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result := struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    []*ChartSet `json:"data"`
	}{
		Data: make([]*ChartSet, 0, amount),
	}
	err = json.Unmarshal(j, &result)
	if err != nil {
		return
	}
	css := result.Data
	s.ChartSets = NewChartSetList(css)
	s.Focus = FocusChartSet
	s.page++
	return
}
