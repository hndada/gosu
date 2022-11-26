package choose

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ModePiano4 = iota
	ModePiano7
	ModeDrum
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

func (s *Scene) LoadChartSetList() (err error) {
	css, err := Search(s.query, s.mode, s.page)
	if err != nil {
		return
	}
	s.ChartSets = NewChartSetList(css)
	s.Focus = FocusChartSet
	s.page++
	return
}
func Search(query string, mode int, page int) (css []*ChartSet, err error) {
	const amount = 1 //20
	const (
		modeMania = 3
		modeTaiko = 1
	)
	u, err := url.Parse("https://api.chimu.moe/v1/search")
	if err != nil {
		return
	}
	vs := url.Values{}
	vs.Add("query", query)
	switch mode {
	case ModePiano4:
		vs.Add("mode", strconv.Itoa(modeMania))
		vs.Add("min_cs", strconv.Itoa(4))
		vs.Add("max_cs", strconv.Itoa(4))
	case ModePiano7:
		vs.Add("mode", strconv.Itoa(modeMania))
		vs.Add("min_cs", strconv.Itoa(7))
		vs.Add("max_cs", strconv.Itoa(7))
	case ModeDrum:
		vs.Add("mode", strconv.Itoa(modeTaiko))
	}
	vs.Add("amount", strconv.Itoa(amount))
	vs.Add("offset", strconv.Itoa(page*amount))
	u.RawQuery = vs.Encode()
	fmt.Println("URL:", u.String())
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
	css = result.Data
	fmt.Println("data length:", len(css))
	return
}
