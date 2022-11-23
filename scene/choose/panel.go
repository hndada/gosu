package choose

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
)

// Duration is in seconds.
// Todo: make UpdateDate time.Time?
type Panel struct {
	bgCh chan draws.Sprite
	draws.Sprite
	// chCh chan mode.Chart
	// mode.Chart

	MusicName  draws.Sprite
	Artist     draws.Sprite
	ChartName  draws.Sprite // Chart-specific
	Charter    draws.Sprite
	UpdateDate draws.Sprite
	Duration   draws.Sprite
	BPM        draws.Sprite
	Level      draws.Sprite // Chart-specific
	// NoteCount  draws.Sprite // Chart-specific

	Inner bool // Enable writing Chart-specific info.
}

func NewPanel(c ChartSet) (p Panel) {
	// Load the image asynchronously.
	p.bgCh = make(chan draws.Sprite)
	go func() {
		i, err := ebitenutil.NewImageFromURL()
		if err != nil {
			return
		}
		s := draws.NewSpriteFromSource(draws.Image{Image: i})
		p.Locate(100, 100, draws.CenterMiddle)
		p.bgCh <- s
		close(p.bgCh)
	}()
	{
		t := draws.NewText(c.Title, Face24)
		p.MusicName = draws.NewSpriteFromSource(t)
	}
	{
		t := draws.NewText(c.Title, Face20)
		p.Artist = draws.NewSpriteFromSource(t)
	}
	{
		t := draws.NewText(c.Creator, Face16)
		p.Charter = draws.NewSpriteFromSource(t)
	}
	{
		if c.RankedStatus >= 1 {
			p.Ranked = true
		}
		t := fmt.Sprintf("last updated at %s", p.UpdateDate)
		if p.Ranked {
			t = fmt.Sprintf("ranked at %s", p.UpdateDate)
		}
	}
	{
		c.ChildrenBeatmaps[0].HitLength
		src := fmt.Sprintf("%02d:%02d", p.Duration/60, p.Duration%60)
		p.Duration = draws.NewSpriteFromSource(t)
	}
	return
}

func (p *Panel) Update(c *Chart) {
	if c != nil {
		if !p.Inner {
			p.updateChart(c)
		}
		p.Inner = true
	} else {
		p.Inner = false
	}
	select {
	case s := <-p.bgCh:
		p.Sprite = s
	default:
	}
}
func (p *Panel) updateChart(c *Chart) {
	{
		t := draws.NewText(c.DiffName, Face20)
		p.ChartName = draws.NewSpriteFromSource(t)
	}
	{ // Todo: use gosu's own level system
		lv := strconv.Itoa(int(c.DifficultyRating * 4))
		t := draws.NewText(fmt.Sprintf("Level: %2d", lv), Face16)
		p.Level = draws.NewSpriteFromSource(t)
	}
	// Due to different logic, MaxCombo tells nothing.
	// Todo: add NoteCount
}
func (p Panel) Draw(dst draws.Image) {
	if p.Sprite.IsValid() {
		p.Sprite.Draw(dst, draws.Op{})
	}
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
