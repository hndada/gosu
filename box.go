package gosu

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

func NewChartBoardBox(c ChartInfo) draws.Box {
	const (
		w = 640
		h = 120
	)
	b := draws.Box{
		Sprite: draws.NewSpriteFromImage(ebiten.NewImage(w, h)),
	}
	{ // Source at left top.
		t := draws.Text{
			Text:  c.MusicSource,
			Color: color.White,
		}
		t.SetFace(Face12, draws.OriginLeftTop)
		b.Texts = append(b.Texts, t)
	}
	{ // Artist at center top.
		t := draws.Text{
			Text:  c.Artist,
			Color: color.White,
		}
		t.SetFace(Face16, draws.OriginCenterMiddle)
		b.Texts = append(b.Texts, t)
	}
	{ // Charter, note counts at right top.
		t := draws.Text{
			Text:  fmt.Sprintf("%s\n◎ %d", c.Charter, c.NoteCounts[0]),
			Color: color.White,
		}
		t.SetFace(Face12, draws.OriginRightTop)
		b.Texts = append(b.Texts, t)
	}
	{ // Music name at center middle.
		t := draws.Text{
			Text:  c.MusicName,
			Color: color.White,
		}
		t.SetFace(Face24, draws.OriginCenterMiddle)
		b.Texts = append(b.Texts, t)
	}
	{ // Time and BPM at left bottom.
		t := draws.Text{
			Text:  fmt.Sprintf("%s\n%s", c.TimeString(), c.BPMString()),
			Color: color.White,
		}
		t.SetFace(Face12, draws.OriginLeftBottom)
		b.Texts = append(b.Texts, t)
	}
	{ // Chart name at center bottom.
		t := draws.Text{
			Text:  c.ChartName,
			Color: color.White,
		}
		t.SetFace(Face20, draws.OriginCenterBottom)
		b.Texts = append(b.Texts, t)
	}
	// Todo: ranked status at right bottom.
	return b
}
func NewChartItemBox(c ChartInfo) draws.Box {
	const (
		w = 500
		h = 80
	)
	// root := draws.Box{
	// 	Sprite:  draws.NewSpriteFromImage(ebiten.NewImage(w, h)),
	// 	PadW:    5,
	// 	PadH:    5,
	// 	MarginW: 5,
	// 	MarginH: 5,
	// }
	root := draws.Box{
		Sprite:  draws.NewSpriteFromImage(ebiten.NewImage(w, h)),
		Pad:     draws.WH{5, 5},
		Margin:  draws.WH{5, 5},
		Content: nil,
	}
	{ // Chart level box.
		t := draws.Text{
			Text:  fmt.Sprintf("%02.f", c.Level),
			Color: color.White,
		}
		t.SetFace(Face20, draws.OriginCenterMiddle)
		b := draws.Box{
			Sprite:  draws.NewSpriteFromImage(ebiten.NewImage(w, h)),
			MarginW: 2,
			MarginH: 2,
			Texts:   []draws.Text{t},
		}
		root.AppendBoxInRow(b)
	}
	{ // Chart text box.
		t := draws.Text{
			Text:  fmt.Sprintf("%s [%s]\n%s / %s", c.MusicName, c.ChartName, c.Artist, c.Charter),
			Color: color.White,
		}
		t.SetFace(Face20, draws.OriginLeftMiddle)
		b := draws.Box{
			Sprite:  draws.NewSpriteFromImage(ebiten.NewImage(w, h)),
			MarginW: 2,
			MarginH: 2,
			Texts:   []draws.Text{t},
		}
		root.AppendBoxInRow(b)
	}
	// Todo: short score box
	return root
}

// Todo: ModeBox, OptionBox, ModsBox, ScoreBox
