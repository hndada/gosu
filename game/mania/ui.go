package mania

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
)

var (
	black = color.RGBA{0, 0, 0, 128}
	red   = color.RGBA{254, 53, 53, 128}
)

type sceneUI struct {
	noteWidths       []int // todo: setNoteSprites()에서만 쓰임
	playfield        game.Sprite
	stageKeys        []game.Sprite
	stageKeysPressed []game.Sprite

	combos      [10]game.Sprite
	scores      [10]game.Sprite
	judgeSprite [len(Judgments)]game.Sprite
	Spotlights  []game.Sprite // 키를 눌렀을 때 불 들어오는 거

	HPBar      game.Sprite // it can go to playfield
	HPBarColor game.Sprite // actually, it can also go to playfield
	HPBarMask  game.Sprite
}

// 가로가 늘어난다고 같이 늘리면 오히려 어색하므로 세로에만 맞춰 늘리기: 100 기준
func newSceneUI(screenSize image.Point, keyCount int) sceneUI {
	s := new(sceneUI)
	scale := float64(screenSize.Y) / 100
	keyKinds := keyKindsMap[keyCount]
	unscaledNoteWidths := Settings.NoteWidths[keyCount&ScratchMask]

	noteWidths := make([]int, keyCount&ScratchMask)
	for key, kind := range keyKinds {
		noteWidths[key] = int(unscaledNoteWidths[kind] * scale)
	}
	i, _ := ebiten.NewImage(screenSize.X, screenSize.Y, ebiten.FilterDefault)

	p := Settings.StagePosition / 100 // proportion
	center := int(float64(screenSize.X) * p)
	var wLeft, wMiddle int
	{ // main
		for _, nw := range noteWidths {
			wMiddle += nw
		}
		h := screenSize.Y
		main, _ := ebiten.NewImage(wMiddle, h, ebiten.FilterDefault)
		main.Fill(color.Black)

		x := center - wMiddle/2 // int - int
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		i.DrawImage(main, op)

		// main := image.NewRGBA(image.Rect(0, 0, w, screenSize.Y))
		// r := image.Rectangle{image.ZP, i.Bounds().Size()}
		// draw.Draw(main, r, &image.Uniform{black}, image.ZP, draw.Over)
	}
	// important: mania-stage-hint에서 판정선이 이미지의 맨 아래에 있다는 보장이 없음
	// 아마 mania-stage-bottom 때문인듯
	// var hHint int
	{ // no-skin ver
		h := int(Settings.JudgeLineHeight * game.DisplayScale())
		hint, _ := ebiten.NewImage(wMiddle, h, ebiten.FilterDefault)
		hint.Fill(red)

		x := center - wMiddle/2 // int - int
		y := int(Settings.HitPosition*game.DisplayScale()) - h
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		i.DrawImage(hint, op)
	}
	// {
	// 	src := Skin.StageHint
	// 	scale := float64(wMiddle) / float64(src.Bounds().Dx())
	// 	h := int(float64(src.Bounds().Dy()) * scale)
	// 	x := center - wMiddle/2
	// 	y := int(Settings.HitPosition*game.DisplayScale()) - h
	// 	op := &ebiten.DrawImageOptions{}
	// 	op.GeoM.Scale(scale, scale)
	// 	op.GeoM.Translate(float64(x), float64(y))
	// 	i.DrawImage(src, op)
	// 	// hHint = h
	// }
	{
		src := Skin.StageLeft
		h := screenSize.Y
		scale := float64(h) / float64(src.Bounds().Dy())
		wLeft = int(float64(src.Bounds().Dx()) * scale)
		x := center - wMiddle/2 - wLeft
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
		i.DrawImage(src, op)
	}
	{
		src := Skin.StageRight
		h := screenSize.Y
		scale := float64(h) / float64(src.Bounds().Dy())
		// wRight = int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
		i.DrawImage(src, op)
	}

	{ // 90도 돌아갈 이미지이므로 whxy 설정에 유의
		src := game.Skin.HPBar
		sprite := game.NewSprite(src)
		// sprite.Theta = 90
		h := int(Settings.HPHeight * game.DisplayScale())
		scale := float64(h) / float64(src.Bounds().Dy())
		w := int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2
		y := game.Settings.ScreenSize.Y - h
		sprite.SetFixedOp(w, h, x, y)
		s.HPBar = sprite
	}
	{ // HP Bar 이미지와 크기가 다를 수 있음
		src := game.Skin.HPBarColor
		sprite := game.NewSprite(src)
		// sprite.Theta = 90
		h := int(Settings.HPHeight * game.DisplayScale())
		scale := float64(h) / float64(src.Bounds().Dy())
		w := int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2 // + s.HPBar.W/2
		y := game.Settings.ScreenSize.Y - h
		// y := int(Settings.HitPosition*game.DisplayScale()) - h
		sprite.SetFixedOp(w, h, x, y)
		s.HPBarColor = sprite

		mask, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
		sprite2 := game.NewSprite(mask)
		sprite2.W = w
		sprite2.H = 0 // hp가 100일 때 0
		sprite2.X = x
		sprite2.Y = y
		sprite2.CompositeMode = ebiten.CompositeModeSourceOut
		s.HPBarMask = sprite2
	}
	s.playfield = game.NewSprite(i)
	s.playfield.SetFixedOp(screenSize.X, screenSize.Y, 0, 0) // todo: 여기에 bg 추가

	s.stageKeys = make([]game.Sprite, keyCount)
	s.stageKeysPressed = make([]game.Sprite, keyCount)

	for k := 0; k < keyCount&ScratchMask; k++ {
		var sprite game.Sprite
		src := Skin.StageKeys[keyKinds[k]]
		sprite = game.NewSprite(src)

		w := noteWidths[k] // 이미지는 크기가 같지만, w가 달라진다

		// scale := float64(sprite.W) / float64(src.Bounds().Size().X)
		// sprite.H = int(float64(src.Bounds().Size().Y) * scale)
		x := center - wMiddle/2 // int - int
		for k2 := 0; k2 < k; k2++ {
			x += noteWidths[k2]
		}
		y := int(Settings.HitPosition * game.DisplayScale()) // + hHint/2
		// fmt.Println(hHint)
		// y := int((Settings.HitPosition - Settings.NoteHeigth/2 -
		// 	4*Settings.NoteHeigth/2) * game.DisplayScale()) // todo: why?
		h := game.Settings.ScreenSize.Y - y

		sprite.SetFixedOp(w, h, x, y)
		s.stageKeys[k] = sprite

		sprite2 := sprite
		src2 := Skin.StageKeysPressed[keyKinds[k]]
		sprite2.SetImage(src2)
		s.stageKeysPressed[k] = sprite2
	}
	{
		src := Skin.StageLight
		sprite := game.NewSprite(src)
		s.Spotlights = make([]game.Sprite, keyCount)
		for k := 0; k < keyCount&ScratchMask; k++ {
			w := noteWidths[k] // 이미지는 크기가 같지만, w가 달라진다
			scale := float64(w) / float64(src.Bounds().Size().X)
			h := int(float64(src.Bounds().Size().Y) * scale)
			x := center - wMiddle/2 // int - int
			for k2 := 0; k2 < k; k2++ {
				x += noteWidths[k2]
			}
			y := int(Settings.HitPosition*game.DisplayScale()) - h
			sprite.Color = Settings.SpotlightColor[keyKinds[k]]
			sprite.SetFixedOp(w, h, x, y)
			s.Spotlights[k] = sprite
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
		x := center - w/2
		y := int(Settings.JudgePosition*game.DisplayScale()) - h/2
		sprite.SetFixedOp(w, h, x, y)
		s.judgeSprite[i] = sprite
	}
	s.noteWidths = noteWidths
	return *s
}

func (s *Scene) setNoteSprites() {
	keyKinds := keyKindsMap[s.chart.KeyCount]

	var wMiddle int
	for k := 0; k < s.chart.KeyCount&ScratchMask; k++ {
		wMiddle += s.noteWidths[k]
	}
	xStart := (game.Settings.ScreenSize.X - wMiddle) / 2
	for i, n := range s.chart.Notes {
		var sprite game.Sprite
		kind := keyKinds[n.Key]
		switch n.Type {
		case TypeNote, TypeLNHead, TypeLNTail: // temp
			sprite = game.NewSprite(Skin.Note[kind])
		}

		scale := float64(s.ScreenSize.Y) / 100
		sprite.H = int(Settings.NoteHeigth * scale)
		sprite.W = s.noteWidths[n.Key]
		x := xStart
		for k := 0; k < n.Key; k++ {
			x += s.noteWidths[k]
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
		ls.Y = tail.Sprite.Y
		ls.Saturation = 1
		ls.Dimness = 1
		s.chart.Notes[i].LongSprite = ls
	}
}
