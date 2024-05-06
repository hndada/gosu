package choose

import (
	"fmt"
	"sort"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

// list: 150x150
// card: 400x140
// cover: 900x250
// slimcover: 1920x360
// Example: https://assets.ppy.sh/beatmaps/784354/covers/slimcover@2x.jpg
// Reference: https://osu.ppy.sh/docs/index.html#beatmapsetcompact-covers
const APIBeatmap = "https://assets.ppy.sh/beatmaps"
const Large = "@2x"

var DefaultCover = draws.NewImage(400, 140)

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

func (c Chart) URLDownload() string {
	return fmt.Sprintf("https://api.chimu.moe/v1/%s", c.DownloadPath)
}

// ChartPanel has own Duration and BPM.
// Todo: chart channel
type ChartPanel struct {
	*ChartSetPanel
	Duration draws.Sprite // in seconds.
	BPM      draws.Sprite

	ChartName draws.Sprite
	Level     draws.Sprite
	NoteCount draws.Sprite
}

func NewChartPanel(sp *ChartSetPanel, c *Chart) *ChartPanel {
	p := &ChartPanel{
		ChartSetPanel: sp,
	}
	{
		second := c.HitLength
		t := fmt.Sprintf("%02d:%02d", second/60, second%60)
		src := draws.NewText(t, scene.Face16)
		s := draws.NewSprite(src)
		s.Locate(450, 0, draws.RightTop)
		p.Duration = s
	}
	{
		bpm := c.BPM
		src := draws.NewText(fmt.Sprintf("%.0f", bpm), scene.Face16)
		s := draws.NewSprite(src)
		s.Locate(450, 50, draws.LeftTop)
		p.BPM = s
	}
	{
		src := draws.NewText(c.DiffName, scene.Face20)
		s := draws.NewSprite(src)
		s.Locate(0, 80, draws.LeftTop)
		p.ChartName = s
	}
	{ // Todo: use gosu's own level system
		lv := Level(c.DifficultyRating)
		src := draws.NewText(fmt.Sprintf("Level: %2d", lv), scene.Face16)
		s := draws.NewSprite(src)
		s.Locate(450, 100, draws.RightTop)
		p.Level = s
	}
	// Todo: NoteCount
	// Due to different logic, MaxCombo tells nothing.
	return p
}

func (p ChartPanel) Draw(dst draws.Image) {
	p.Sprite.Draw(dst, draws.Op{})
	p.MusicName.Draw(dst, draws.Op{})
	p.Artist.Draw(dst, draws.Op{})
	p.ChartName.Draw(dst, draws.Op{})
	p.Charter.Draw(dst, draws.Op{})
	p.UpdateDate.Draw(dst, draws.Op{})

	p.Duration.Draw(dst, draws.Op{})
	p.BPM.Draw(dst, draws.Op{})
	p.Level.Draw(dst, draws.Op{})
	// p.NoteCount.Draw(dst, draws.Op{})
}
