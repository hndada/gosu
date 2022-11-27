package choose

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	const rankedOnly = true
	css, err := Search(s.query, s.mode, s.page, s.levelLimit, rankedOnly)
	if err != nil {
		return
	}
	s.ChartSets = NewChartSetList(css)
	s.Focus = FocusChartSet
	return
}
func Search(query string, mode, page int, lvLimit, rankedOnly bool) (css []*ChartSet, err error) {
	const amount = RowCount // = 20
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
	if lvLimit {
		vs.Add("max_diff", "4")
	}
	if rankedOnly {
		vs.Add("status", "1,2,3")
	}
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
	if len(css) < amount {
		css2, _ := Search(query, mode, page, lvLimit, false)
		css = append(css, css2...)
		css = css[:amount]
	}
	// fmt.Println("data length:", len(css))
	return
}
