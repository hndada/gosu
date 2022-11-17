package choose

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/db"
)

// Todo: different face on same Sprite (or Box, or sth)
func NewChartPanel(c db.Chart) (panel draws.Sprite) {
	i := draws.NewImage(800, 225)
	i.Fill(color.Black)
	panel = draws.NewSpriteFromSource(i)
	panel.Locate(100, 100, draws.CenterMiddle)
	{
		t := c.MusicName
		text := draws.NewText(t, draws.LoadDefaultFace(24))
		src := draws.NewSpriteFromSource(text)
		panel.Append(src, draws.Location{X: 0.01, Y: 0.05})
	}
	{
		t := c.Artist
		text := draws.NewText(t, draws.LoadDefaultFace(16))
		src := draws.NewSpriteFromSource(text)
		panel.Append(src, draws.Location{X: 0.01, Y: 0.25})
	}
	{
		t := c.ChartName
		text := draws.NewText(t, draws.LoadDefaultFace(18))
		src := draws.NewSpriteFromSource(text)
		panel.Append(src, draws.Location{X: 0.01, Y: 0.65})
	}
	{
		t := c.Charter
		text := draws.NewText(t, draws.LoadDefaultFace(14))
		src := draws.NewSpriteFromSource(text)
		panel.Append(src, draws.Location{X: 0.01, Y: 0.85})
	}
	{
		var (
			status     string
			ranked     bool
			updateTime time.Timer
		)
		if ranked {
			status = "ranked at %v"
		} else {
			status = "last updated at %v"
		}
		t := fmt.Sprintf(status, updateTime)
		text := draws.NewText(t, draws.LoadDefaultFace(12))
		src := draws.NewSpriteFromSource(text)
		panel.Append(src, draws.Location{X: 0.01, Y: 1})
	}
	{
		minute, second := c.Duration/60000, c.Duration%60000/1000
		t := fmt.Sprintf("%02d:%02d\n%.0f BPM\n(%.0f-%.0f)",
			minute, second,
			c.MainBPM,
			c.MinBPM, c.MaxBPM,
		)
		text := draws.NewText(t, draws.LoadDefaultFace(14))
		src := draws.NewSpriteFromSource(text)
		panel.Append(src, draws.Location{X: 0.85, Y: 0.05})
	}
	{
		t := fmt.Sprintf("Level: %4.2f\nNotes count: %d",
			c.Level,
			c.NoteCounts[0],
		)
		text := draws.NewText(t, draws.LoadDefaultFace(14))
		src := draws.NewSpriteFromSource(text)
		panel.Append(src, draws.Location{X: 0.85, Y: 1})
	}
	fmt.Printf("%+v\n", panel)
	return
}
