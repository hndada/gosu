package choose

import (
	_ "embed"
	"fmt"
	"sort"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

type ChartSet struct {
	SetId            int
	ChildrenBeatmaps []*Chart
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

//go:embed proxy.txt
var proxy string

func (c ChartSet) URLCover(kind, suffix string) string {
	return fmt.Sprintf("%s%s/%d/covers/%s%s.jpg", proxy, APIBeatmap, c.SetId, kind, suffix) // https://proxy.cors.sh/
}
func (c ChartSet) URLPreview() string {
	return fmt.Sprintf("%shttps://b.ppy.sh/preview/%d.mp3", proxy, c.SetId) // https://proxy.cors.sh/
}
func (c ChartSet) URLDownload() string {
	return fmt.Sprintf("%shttps://api.chimu.moe/v1/d/%d", proxy, c.SetId) // https://proxy.cors.sh/
}

type ChartSetList struct {
	*List
	ChartSets []*ChartSet
	// Panel     *ChartSetPanel
}

func NewChartSetList(css []*ChartSet) (l ChartSetList) {
	rows := make([]Row, len(css))
	for i, cs := range css {
		card := cs.URLCover("card", "")
		thumb := cs.URLCover("list", "")
		rows[i] = NewRow(card, thumb, cs.Title, cs.Artist)
	}
	sort.Slice(rows, func(i, j int) bool {
		return css[i].LastUpdate > css[j].LastUpdate
	})
	l.List = NewList(rows)
	l.ChartSets = css
	// l.Panel = NewChartSetPanel(l.Current())
	return
}
func (l *ChartSetList) Update() (fired bool, state int) {
	fired, state = l.List.Update()
	// if l.Panel != nil {
	// 	l.Panel.Update()
	// }
	// if fired = l.Cursor.Update(); fired {
	// 	cs := l.ChartSets[l.cursor]
	// 	l.Panel = NewChartSetPanel(cs)
	// }
	return
}
func (l ChartSetList) Draw(dst draws.Image) {
	l.List.Draw(dst)
	// l.Panel.Draw(dst)
}
func (l ChartSetList) Current() *ChartSet {
	if len(l.ChartSets) == 0 {
		return nil
	}
	return l.ChartSets[l.cursor]
}

type ChartSetPanel struct {
	// bgCh chan draws.Image
	draws.Sprite

	MusicName  draws.Sprite
	Artist     draws.Sprite
	Charter    draws.Sprite
	UpdateDate draws.Sprite
	Duration   draws.Sprite // in seconds.
	BPM        draws.Sprite
}

var DefaultCover = draws.NewImage(400, 140)

func NewChartSetPanel(cs *ChartSet) *ChartSetPanel {
	p := &ChartSetPanel{}
	// p.bgCh = make(chan draws.Image)
	{
		s := draws.NewSprite(DefaultCover)
		s.Locate(100, 100, draws.CenterMiddle)
		p.Sprite = s
	}
	go func() {
		i, err := draws.NewImageFromURL(cs.URLCover("cover", Large))
		if err != nil {
			fmt.Println("chart set cover: ", err)
		}
		s := draws.NewSprite(i)
		p.Sprite.Source = s
		// p.bgCh <- draws.Image{Image: i}
		// close(p.bgCh)
	}()
	{
		src := draws.NewText(cs.Title, scene.Face24)
		s := draws.NewSprite(src)
		s.Locate(0, 0, draws.LeftTop)
		p.MusicName = s
	}
	{
		src := draws.NewText(cs.Artist, scene.Face20)
		s := draws.NewSprite(src)
		s.Locate(0, 25, draws.LeftTop)
		p.Artist = s
	}
	{
		src := draws.NewText(cs.Creator, scene.Face16)
		s := draws.NewSprite(src)
		s.Locate(0, 100, draws.LeftTop)
		p.Charter = s
	}
	{
		format := "Last updated at %s"
		if cs.RankedStatus >= 1 {
			format = "ranked at %s"
		}
		src := draws.NewText(fmt.Sprintf(format, p.UpdateDate), scene.Face16)
		s := draws.NewSprite(src)
		s.Locate(0, 125, draws.LeftTop)
		p.UpdateDate = s
	}
	{
		second := cs.ChildrenBeatmaps[0].HitLength
		t := fmt.Sprintf("%02d:%02d", second/60, second%60)
		src := draws.NewText(t, scene.Face12)
		s := draws.NewSprite(src)
		s.Locate(450, 0, draws.RightTop)
		p.Duration = s
	}
	{
		bpm := cs.ChildrenBeatmaps[0].BPM
		src := draws.NewText(fmt.Sprintf("%.0f", bpm), scene.Face12)
		s := draws.NewSprite(src)
		s.Locate(450, 50, draws.LeftTop)
		p.BPM = s
	}
	return p
}
func (p *ChartSetPanel) Update() {
	// select {
	// case i := <-p.bgCh:
	// 	s := draws.NewSprite(i)
	// 	s.Locate(100, 100, draws.CenterMiddle)
	// 	p.Sprite = s
	// default:
	// }
}
func (p ChartSetPanel) Draw(dst draws.Image) {
	p.Sprite.Draw(dst, draws.Op{})
	p.MusicName.Draw(dst, draws.Op{})
	p.Artist.Draw(dst, draws.Op{})
	p.Charter.Draw(dst, draws.Op{})
	p.UpdateDate.Draw(dst, draws.Op{})
	p.Duration.Draw(dst, draws.Op{})
	p.BPM.Draw(dst, draws.Op{})
}
