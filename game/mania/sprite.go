package mania

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
)

var (
	black = color.RGBA{0, 0, 0, 128}
	red   = color.RGBA{254, 53, 53, 128}
)

type sceneUI struct {
	noteWidths []int
	playfield  game.Sprite
}

// 가로가 늘어난다고 같이 늘리면 오히려 어색하므로 세로에만 맞춰 늘리기: 100 기준
func newSceneUI(screenSize image.Point, keyCount int) sceneUI {
	s := new(sceneUI)
	scale := float64(screenSize.Y) / 100
	keyKinds := keyKindsMap[keyCount]
	unscaledNoteWidths := Settings.NoteWidths[keyCount&ScratchMask]

	s.noteWidths = make([]int, keyCount&ScratchMask)
	for key, kind := range keyKinds {
		s.noteWidths[key] = int(unscaledNoteWidths[kind] * scale)
	}

	{ // playfield
		var w int
		for _, nw := range s.noteWidths {
			w += nw
		}

		i := image.NewRGBA(image.Rect(0, 0, w, screenSize.Y))
		{ // main
			r := image.Rectangle{image.ZP, i.Bounds().Size()}
			draw.Draw(i, r, &image.Uniform{black}, image.ZP, draw.Over)
		}
		{ // hint
			hp := int(Settings.HitPosition * scale)
			h := int(Settings.NoteHeigth * scale)
			sp := image.Point{0, hp - h/2}
			r := image.Rectangle{sp, sp.Add(image.Pt(w, h))}
			draw.Draw(i, r, &image.Uniform{red}, image.ZP, draw.Over)
		}
		s.playfield.SetImage(i)
		p := Settings.StagePosition / 100                  // position in proportion
		s.playfield.X = int(float64(screenSize.X)*p) - w/2 // int - int
		s.playfield.Y = 0
		s.playfield.W = w
		s.playfield.H = screenSize.Y
	}
	return *s
}

func (s sceneUI) Draw(screen *ebiten.Image) {
	s.playfield.Draw(screen)
}

// todo: n.scored 시 명암 처리
// length: 시간적 길이 // height: 공간적 길이. 이미지 길이.
func (s *Scene) setNoteSprites() {
	var sprite game.Sprite
	keyKinds := keyKindsMap[s.chart.KeyCount]
	for i, n := range s.chart.Notes {
		kind := keyKinds[n.Key]
		switch n.Type {
		case TypeNote, TypeLNHead, TypeLNTail: // temp
			sprite.SetImage(Skin.Note[kind])
		}

		scale := float64(s.ScreenSize.Y) / 100
		sprite.H = int(Settings.NoteHeigth * scale)
		sprite.W = s.ui.noteWidths[n.Key]
		x := s.ui.playfield.X
		for k := 0; k < n.Key; k++ {
			x += s.ui.noteWidths[k]
		}
		sprite.X = x
		y := Settings.HitPosition - n.position*s.speed - float64(sprite.H)/2
		sprite.Y = int(y * scale)
		s.chart.Notes[i].Sprite = sprite
	}

	// LN body sprite
	// 모든 Sprite는 자신의 값을 갱신 시켜줄 개체와 connect되어 있어야 함
	kinds := keyKindsMap[s.chart.KeyCount]
	for i, tail := range s.chart.Notes {
		if tail.Type != TypeLNTail {
			continue
		}
		head := s.chart.Notes[tail.prev]
		ls := game.LongSprite{
			Vertical: true,
		}
		ls.SetImage(Skin.LNBody[kinds[tail.Key]]) // temp: no animation support
		ls.W = tail.Sprite.W
		ls.H = head.Sprite.Y - tail.Sprite.Y
		ls.X = tail.Sprite.X
		ls.Y = tail.Sprite.Y // + 50*tail.Sprite.H // todo: why?
		s.chart.Notes[i].LongSprite = ls
	}
}
