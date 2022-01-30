package common

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/engine/ui"
)

type Judgment struct {
	Value   float64
	Penalty float64
	HP      float64
	Window  int64
}

type JudgmentMeter struct {
	Judgments []Judgment
	Sprite    ui.FixedSprite
}

var (
	brown     = color.RGBA{156, 42, 42, 255}
	yellow    = color.RGBA{184, 134, 11, 255}
	green     = color.RGBA{0, 255, 0, 255}
	lightblue = color.RGBA{0, 181, 204, 255}
	blue      = color.RGBA{0, 0, 255, 255}
)

// temp
var judgmentMeterColor []color.RGBA = []color.RGBA{brown, yellow, green, lightblue, blue}

const JudgmentMeterScale = 2

func NewJudgmentMeter(js []Judgment) *JudgmentMeter {
	jm := new(JudgmentMeter)
	jm.Judgments = js

	var s ui.Sprite
	var base *ebiten.Image
	{ // set base box
		// TODO: 검은색 바탕 상자가 안 그려진다
		const height = 5 // 높이는 세로 전체 100 기준 5
		j := jm.Judgments[len(jm.Judgments)-1]
		w := int(JudgmentMeterScale*float64(j.Window)) * 2
		h := int(DisplayScale() * height)

		base = ebiten.NewImage(w, h)
		base.Fill(color.RGBA64{0, 0, 0, 255})

		s = ui.NewSprite(base) // base is just for providing size info
		s.W = w
		s.H = h
		s.X = Settings.ScreenSizeX/2 - s.W/2
		s.Y = Settings.ScreenSizeY - s.H
	}
	{ // set color box
		const height = 1 // base 대비 1
		h := int(DisplayScale() * height)
		y := (base.Bounds().Dy() - h) / 2
		for i := range jm.Judgments {
			j := jm.Judgments[len(jm.Judgments)-1-i]
			w := int(JudgmentMeterScale*float64(j.Window)) * 2
			x := base.Bounds().Dx()/2 - w/2

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			box := ebiten.NewImage(w, h)
			box.Fill(judgmentMeterColor[i])
			base.DrawImage(box, op)
		}
	}
	{ // set middle line
		const height = 5 // 높이는 세로 전체 100 기준 5
		w := 1
		h := int(DisplayScale() * height)
		x := base.Bounds().Dx()/2 - w
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		box := ebiten.NewImage(w, h)
		box.Fill(color.White)
		base.DrawImage(box, op)
	}
	s.SetImage(base)
	jm.Sprite = ui.NewFixedSprite(s)
	return jm
}

// // "early" goes plus
// // TODO: 종종 x값이 음수가 나옴. 저 멀리의 노트로 timeDiff를 계산하는 걸수도 있음
// func (jm JudgmentMeter) NewTimingSprite(timeDiff int64) ui.Animation {
// 	w := int(Settings.JudgmentMeterScale)
// 	h := jm.Sprite.H
// 	x := Settings.ScreenSizeX/2 - int(Settings.JudgmentMeterScale*float64(timeDiff))
// 	y := jm.Sprite.Y

// 	i := image.NewRGBA(image.Rect(0, 0, w, h))
// 	r := image.Rectangle{image.ZP, i.Bounds().Size()}
// 	draw.Draw(i, r, &image.Uniform{color.RGBA{255, 255, 255, 128}}, image.ZP, draw.Over)
// 	i := ebiten.NewImageFromImage(i)
// 	// i := ebiten.NewImage(w, h)
// 	// i.Fill(color.White)

// 	a := ui.NewAnimation([]*ebiten.Image{i})
// 	a.W = w
// 	a.H = h
// 	a.X = x
// 	a.Y = y
// 	a.Rep = 20 // temp
// 	return a
// }
