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

	var sprite Sprite
	var base *ebiten.Image

	{ // set base box
		// todo: 검은색 바탕 상자가 안 그려진다
		const height = 5 // 높이는 세로 전체 100 기준 5
		j := jm.Judgments[len(jm.Judgments)-1]
		w := int(Settings.JudgmentMeterScale*float64(j.Window)) * 2
		h := int(DisplayScale() * height)
		x := Settings.ScreenSize.X/2 - w/2
		y := Settings.ScreenSize.Y - h

		base, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		base.Fill(color.RGBA64{0, 0, 0, 255})
		sprite = NewSprite(base) // base is just for providingsize info
		sprite.SetFixedOp(w, h, x, y)
	}
	{ // set color box
		const height = 1 // base 대비 1
		h := int(DisplayScale() * height)
		y := (base.Bounds().Dy() - h) / 2
		for i := range jm.Judgments {
			j := jm.Judgments[len(jm.Judgments)-1-i]
			w := int(Settings.JudgmentMeterScale*float64(j.Window)) * 2
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
		h := int(DisplayScale() * height)
		x := base.Bounds().Dx()/2 - w
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		box, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
		box.Fill(color.White)
		_ = base.DrawImage(box, op)
	}
	sprite.SetImage(base)
	jm.Sprite = sprite
	return jm
}

// "early" goes plus
// todo: 종종 x값이 음수가 나옴. 저 멀리의 노트로 timeDiff를 계산하는 걸수도 있음
func (jm JudgmentMeter) NewTimingSprite(timeDiff int64) Sprite {
	w := int(Settings.JudgmentMeterScale)
	h := jm.Sprite.H
	i, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	i.Fill(color.White)
	x := Settings.ScreenSize.X/2 - int(Settings.JudgmentMeterScale*float64(timeDiff))
	y := jm.Sprite.Y
	sprite := NewSprite(i)
	sprite.SetFixedOp(w, h, x, y)
	// fmt.Println(w, h, x, y)
	return sprite
}
