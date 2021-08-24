package mania

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hndada/gosu/game"
)

var (
	black = color.RGBA{0, 0, 0, 128}
	red   = color.RGBA{254, 53, 53, 128}
)

type sceneUI struct {
	noteWidths       []int
	playfield        game.Sprite
	stageKeys        []game.Sprite
	stageKeysPressed []game.Sprite

	combos      [10]game.Sprite
	scores      [10]game.Sprite
	judgeSprite [len(Judgments)]game.Sprite
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
		p := Settings.StagePosition / 100 // position in proportion
		s.playfield.W = w
		s.playfield.H = screenSize.Y
		s.playfield.X = int(float64(screenSize.X)*p) - w/2 // int - int
		s.playfield.Y = 0
		s.playfield.Saturation = 1
		s.playfield.Dimness = 1
	}
	{
		s.stageKeys = make([]game.Sprite, keyCount)
		s.stageKeysPressed = make([]game.Sprite, keyCount)
		for k := 0; k < keyCount&ScratchMask; k++ {
			var sprite game.Sprite
			src := Skin.StageKeys[keyKinds[k]]
			sprite = game.NewSprite(src)

			w := s.noteWidths[k] // 이미지는 크기가 같지만, w가 달라진다

			// scale := float64(sprite.W) / float64(src.Bounds().Size().X)
			// sprite.H = int(float64(src.Bounds().Size().Y) * scale)
			y := int((Settings.HitPosition - Settings.NoteHeigth/2 -
				4*Settings.NoteHeigth/2) * game.DisplayScale()) // todo: why?
			h := game.Settings.ScreenSize.Y - y
			x := s.playfield.X
			for k2 := 0; k2 < k; k2++ {
				x += s.noteWidths[k2]
			}
			sprite.SetFixedOp(w, h, x, y)
			s.stageKeys[k] = sprite

			sprite2 := sprite
			src2 := Skin.StageKeysPressed[keyKinds[k]]
			sprite2.SetImage(src2)
			s.stageKeysPressed[k] = sprite2
		}
	}

	s.combos = game.LoadNumbers(game.NumberCombo)
	s.scores = game.LoadNumbers(game.NumberScore)

	for i := range s.judgeSprite {
		src := Skin.Judge[i]
		sprite := game.NewSprite(src)
		h := int(Settings.JudgeHeight * game.DisplayScale())
		scale := float64(h) / float64(src.Bounds().Dy())
		w := int(float64(src.Bounds().Dx()) * scale)
		x := (game.Settings.ScreenSize.X - w) / 2
		y := int(Settings.JudgePosition*game.DisplayScale()) - h/2
		sprite.SetFixedOp(w, h, x, y)
		s.judgeSprite[i] = sprite
	}
	return *s
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
		sprite.W = s.noteWidths[n.Key]
		x := s.playfield.X
		for k := 0; k < n.Key; k++ {
			x += s.noteWidths[k]
		}
		sprite.X = x
		y := Settings.HitPosition - n.position*s.speed - float64(sprite.H)/2
		sprite.Y = int(y * scale)
		sprite.Saturation = 1
		sprite.Dimness = 1
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
		ls.Y = tail.Sprite.Y
		ls.Saturation = 1
		ls.Dimness = 1
		s.chart.Notes[i].LongSprite = ls
	}
}
