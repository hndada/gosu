package choose

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

type ChartSetPanel struct {
	bgCh chan draws.Sprite
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
	// Load the image asynchronously.
	p.bgCh = make(chan draws.Sprite)
	go func() {
		i, err := ebitenutil.NewImageFromURL(cs.URLCover("cover", Large))
		if err != nil {
			return
		}
		s := draws.NewSpriteFromSource(draws.Image{Image: i})
		p.Locate(100, 100, draws.CenterMiddle)
		p.bgCh <- s
		close(p.bgCh)
	}()
	{
		src := draws.NewText(cs.Title, scene.Face24)
		s := draws.NewSpriteFromSource(src)
		s.Locate(0, 0, draws.LeftTop)
		p.MusicName = s
	}
	{
		src := draws.NewText(cs.Artist, scene.Face20)
		s := draws.NewSpriteFromSource(src)
		s.Locate(0, 25, draws.LeftTop)
		p.Artist = s
	}
	{
		src := draws.NewText(cs.Creator, scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(0, 100, draws.LeftTop)
		p.Charter = s
	}
	{
		format := "Last updated at %s"
		if cs.RankedStatus >= 1 {
			format = "ranked at %s"
		}
		src := draws.NewText(fmt.Sprintf(format, p.UpdateDate), scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(0, 125, draws.LeftTop)
		p.UpdateDate = s
	}
	{
		second := cs.ChildrenBeatmaps[0].HitLength
		t := fmt.Sprintf("%02d:%02d", second/60, second%60)
		src := draws.NewText(t, scene.Face12)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 0, draws.RightTop)
		p.Duration = s
	}
	{
		bpm := cs.ChildrenBeatmaps[0].BPM
		src := draws.NewText(fmt.Sprintf("%.0f", bpm), scene.Face12)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 50, draws.LeftTop)
		p.BPM = s
	}
	return p
}
func (p *ChartSetPanel) Update() {
	select {
	case s := <-p.bgCh:
		p.Sprite = s
	default:
	}
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

func NewChartPanel(csp *ChartSetPanel, c *Chart) *ChartPanel {
	p := &ChartPanel{
		ChartSetPanel: csp,
	}
	{
		second := c.HitLength
		t := fmt.Sprintf("%02d:%02d", second/60, second%60)
		src := draws.NewText(t, scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 0, draws.RightTop)
		p.Duration = s
	}
	{
		bpm := c.BPM
		src := draws.NewText(fmt.Sprintf("%.0f", bpm), scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 50, draws.LeftTop)
		p.BPM = s
	}
	{
		src := draws.NewText(c.DiffName, scene.Face20)
		s := draws.NewSpriteFromSource(src)
		s.Locate(0, 80, draws.LeftTop)
		p.ChartName = s
	}
	{ // Todo: use gosu's own level system
		lv := strconv.Itoa(int(c.DifficultyRating * 4))
		src := draws.NewText(fmt.Sprintf("Level: %2d", lv), scene.Face16)
		s := draws.NewSpriteFromSource(src)
		s.Locate(450, 100, draws.RightTop)
		p.Level = s
	}
	// Todo: NoteCount
	// Due to different logic, MaxCombo tells nothing.
	return p
}
func (p *ChartPanel) Update(cs *ChartSet, c *Chart) {
	p.ChartSetPanel.Update()
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
