package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type Judgment struct {
	Value   float64
	Penalty float64
	HP      float64
	Window  int64
}

type JudgmentMeter struct {
	Judgments []Judgment
	Sprite    Sprite
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

func NewJudgmentMeter(js []Judgment) *JudgmentMeter {
	jm := new(JudgmentMeter)
	jm.Judgments = js

	const scale = 2 // 1ms 당 2px
	var sprite Sprite
	var base *ebiten.Image

	{ // set base box
		// todo: 검은색 바탕 상자가 안 그려진다
		const height = 5 // 높이는 세로 전체 100 기준 5
		j := jm.Judgments[len(jm.Judgments)-1]
		w := int(scale*j.Window) * 2
		h := int(Scale() * height)
		x := Settings.ScreenSize.X/2 - w/2
		y := Settings.ScreenSize.Y - h
		sprite.Op = &ebiten.DrawImageOptions{}
		sprite.Op.GeoM.Translate(float64(x), float64(y))
		sprite.W = w // todo: WHXY 과 Op 한번에 결정
		sprite.H = h
		sprite.X = x
		sprite.Y = y
		base, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		base.Fill(color.RGBA64{0, 0, 0, 255})
	}
	{ // set color box
		const height = 1 // base 대비 1
		h := int(Scale() * height)
		y := (base.Bounds().Dy() - h) / 2
		for i := range jm.Judgments {
			j := jm.Judgments[len(jm.Judgments)-1-i]
			w := int(scale*j.Window) * 2
			x := base.Bounds().Dx()/2 - w/2

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			box, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
			box.Fill(judgmentMeterColor[i])
			_ = base.DrawImage(box, op)
		}
	}
	{ // set middle line
		const height = 5 // 높이는 세로 전체 100 기준 5
		w := 1
		h := int(Scale() * height)
		x := base.Bounds().Dx()/2 - w
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		box, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
		box.Fill(color.White)
		_ = base.DrawImage(box, op)
	}
	sprite.SetEbitenImage(base)
	jm.Sprite = sprite
	return jm
}

// "early" goes plus
func (jm JudgmentMeter) DrawTiming(screen *ebiten.Image, timeDiffs []int64) {
	const scale = 2 // 1ms 당 2px
	for _, t := range timeDiffs {
		w := scale
		h := jm.Sprite.H
		i, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
		i.Fill(color.White)

		x := Settings.ScreenSize.X/2 - int(scale*t)
		y := jm.Sprite.Y
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(i, op)
	}
}
