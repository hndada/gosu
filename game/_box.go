package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/framework/draws"
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

func (c ChartInfo) NewChartBoard() draws.Box {
	var (
		ws = []float64{140, 360, 140}
		hs = []float64{24, 36, 32, 24}
	)
	boxs := [][]draws.Box{
		{
			{
				Inner: draws.NewLabel(c.MusicSource, Face12, color.White),
				Align: draws.AtMin,
			},
			{
				Inner: draws.NewLabel(c.Artist, Face12, color.White),
				Align: draws.ModeXY{X: draws.ModeMid, Y: draws.ModeMin},
			},
			{
				Inner: draws.NewLabel(c.NoteCountString(), Face12, color.White),
				Align: draws.ModeXY{X: draws.ModeMax, Y: draws.ModeMin},
			},
		},
		{
			{
				// Inner: draws.NewRectangle(draws.Pt(ws[0], hs[1])),
			},
			{
				Inner: draws.NewLabel(c.MusicName, Face20, color.White),
				Align: draws.ModeXY{X: draws.ModeMid, Y: draws.ModeMax},
			},
			{
				// Inner: draws.NewRectangle(draws.Pt(ws[2], hs[1])),
			},
		},
		{
			{
				Inner: draws.NewLabel(c.TimeString(), Face16, color.White),
				Align: draws.AtMin,
			},
			{
				Inner: draws.NewLabel(c.ChartName, Face16, color.White),
				Align: draws.ModeXY{X: draws.ModeMid, Y: draws.ModeMin},
			},
			{
				// Inner: draws.NewRectangle(draws.Pt(ws[2], hs[2])),
			},
		},
		{
			{
				Inner: draws.NewLabel(c.BPMString(), Face16, color.White),
				Align: draws.ModeXY{X: draws.ModeMin, Y: draws.ModeMax},
			},
			{
				Inner: draws.NewLabel(c.Charter, Face12, color.White),
				Align: draws.ModeXY{X: draws.ModeMid, Y: draws.ModeMax},
			},
			{ // Todo: ranked status
				// Inner: draws.NewRectangle(draws.Pt(ws[2], hs[3])),
				Align: draws.ModeXY{X: draws.ModeMax, Y: draws.ModeMax},
			},
		},
	}
	for i, row := range boxs {
		for j := range row {
			boxs[i][j].Outer = draws.NewRectangle(draws.Pt(ws[j], hs[i]))
			// fmt.Printf("%d %d: %+v\n", i, j, box)
		}
		// fmt.Println()
	}
	// boxs = make([][]draws.Box, 3)
	// for i := range boxs {
	// 	boxs[i] = make([]draws.Box, 3)
	// 	for j := range boxs[i] {
	// 		boxs[i][j] = draws.Box{
	// 			Inner: labels[i][j],
	// 			Align: draws.ModeXY{i, j},
	// 		}
	// 	}
	// }

	board := draws.Box{
		Inner:   draws.NewGrid(boxs, ws, hs, draws.Point{}),
		Pad:     draws.Pt(10, 10),
		Point:   draws.Pt(ScreenSizeX/2, 150),
		Origin2: draws.AtMid,
		Align:   draws.AtMid,
	}
	board.Outer = &draws.Rectangle{
		Size_: board.OuterSize(),
		Color: color.NRGBA{128, 128, 128, 128},
	}
	// box.Outer = draws.NewRectangle(box.OuterSize(), gray)
	// outerImage := ebiten.NewImage(box.OuterSize().XYInt())
	// outerImage.Fill(gray)
	// box.Outer = draws.NewSprite3FromImage(outerImage)
	return board
}
