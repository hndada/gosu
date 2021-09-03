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

type TimeBool struct {
	Time  int64
	Value bool
}
type sceneUI struct {
	noteWidths       []int // todo: setNoteSprites()에서만 쓰임
	playfield        game.FixedSprite
	stageKeys        []game.FixedSprite
	stageKeysPressed []game.FixedSprite

	combos      [10]game.Sprite
	scores      [10]game.Sprite
	judgeSprite [len(Judgments)]game.Animation // todo: rename
	Spotlights  []game.FixedSprite             // 키를 눌렀을 때 불 들어오는 거

	HPBar      game.FixedSprite // it can be in playfield
	HPBarColor game.FixedSprite // actually, it can also go to playfield
	HPBarMask  game.Sprite
	hpScreen   *ebiten.Image

	Lighting   []game.Animation // 여러 lane에서 동시에 그려져야함
	LightingLN []game.Animation

	jm            *game.JudgmentMeter // temp
	timingSprites []game.Animation    // temp

	bg game.FixedSprite
}

// 가로가 늘어난다고 같이 늘리면 오히려 어색하므로 세로에만 맞춰 늘리기: 100 기준
func newSceneUI(c *Chart, keyCountWithScratchMode int) sceneUI {
	s := new(sceneUI)
	keyCount := keyCountWithScratchMode & ScratchMask // temp
	scale := float64(game.Settings.ScreenSize.Y) / 100
	keyKinds := keyKindsMap[keyCount]
	unscaledNoteWidths := Settings.NoteWidths[keyCount]

	s.bg = c.BG(game.Settings.BackgroundDimness)

	noteWidths := make([]int, keyCount)
	for key, kind := range keyKinds {
		noteWidths[key] = int(unscaledNoteWidths[kind] * scale)
	}
	i, _ := ebiten.NewImage(game.Settings.ScreenSize.X, game.Settings.ScreenSize.Y, ebiten.FilterDefault)

	p := Settings.StagePosition / 100 // proportion
	center := int(float64(game.Settings.ScreenSize.X) * p)
	var wLeft, wMiddle int
	{ // main
		for _, nw := range noteWidths {
			wMiddle += nw
		}
		h := game.Settings.ScreenSize.Y

		// temp: Fill이 alpha value를 안 받는 것 같아 draw.Draw 사용 중
		mainSrc := image.NewRGBA(image.Rect(0, 0, wMiddle, h))
		r := image.Rectangle{image.ZP, i.Bounds().Size()}
		draw.Draw(mainSrc, r, &image.Uniform{black}, image.ZP, draw.Over)
		main, _ := ebiten.NewImageFromImage(mainSrc, ebiten.FilterDefault)
		// main, _ := ebiten.NewImage(wMiddle, h, ebiten.FilterDefault)
		// main.Fill(color.Black)

		x := center - wMiddle/2 // int - int
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.ColorM.Scale(0, 0, 0, 1)

		const dimness = 30 // temp
		op.ColorM.ChangeHSV(0, 1, float64(dimness)/100)
		i.DrawImage(main, op)
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
		h := game.Settings.ScreenSize.Y
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
		h := game.Settings.ScreenSize.Y
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
		sprite := game.NewFixedSprite(src)
		// sprite.Theta = 90
		h := int(Settings.HPHeight * game.DisplayScale())
		scale := float64(h) / float64(src.Bounds().Dy())
		w := int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2
		y := game.Settings.ScreenSize.Y - h
		sprite.W = w
		sprite.H = h
		sprite.X = x
		sprite.Y = y
		sprite.Fix()
		s.HPBar = sprite
	}
	{ // HP Bar 이미지와 크기가 다를 수 있음
		src := game.Skin.HPBarColor
		sprite := game.NewFixedSprite(src)
		// sprite.Theta = 90
		h := int(Settings.HPHeight * game.DisplayScale())
		scale := float64(h) / float64(src.Bounds().Dy())
		w := int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2 // + s.HPBar.W/2
		y := game.Settings.ScreenSize.Y - h
		// y := int(Settings.HitPosition*game.DisplayScale()) - h
		sprite.W = w
		sprite.H = h
		sprite.X = x
		sprite.Y = y
		sprite.Fix()
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
	s.playfield = game.NewFixedSprite(i)
	s.playfield.W = game.Settings.ScreenSize.X
	s.playfield.H = game.Settings.ScreenSize.Y
	s.playfield.X = 0
	s.playfield.Y = 0
	s.playfield.Fix() // todo: 여기에 bg 추가

	s.stageKeys = make([]game.FixedSprite, keyCount)
	s.stageKeysPressed = make([]game.FixedSprite, keyCount)

	for k := 0; k < keyCount; k++ {
		var sprite game.FixedSprite
		src := Skin.StageKeys[keyKinds[k]]
		sprite = game.NewFixedSprite(src)

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

		sprite.W = w
		sprite.H = h
		sprite.X = x
		sprite.Y = y
		sprite.Fix()
		s.stageKeys[k] = sprite

		sprite2 := sprite
		src2 := Skin.StageKeysPressed[keyKinds[k]]
		sprite2.SetImage(src2)
		s.stageKeysPressed[k] = sprite2
	}
	{
		src := Skin.StageLight
		sprite := game.NewFixedSprite(src)
		s.Spotlights = make([]game.FixedSprite, keyCount)
		for k := 0; k < keyCount; k++ {
			w := noteWidths[k] // 이미지는 크기가 같지만, w가 달라진다
			scale := float64(w) / float64(src.Bounds().Size().X)
			h := int(float64(src.Bounds().Size().Y) * scale)
			x := center - wMiddle/2 // int - int
			for k2 := 0; k2 < k; k2++ {
				x += noteWidths[k2]
			}
			y := int(Settings.HitPosition*game.DisplayScale()) - h
			sprite.Color = Settings.SpotlightColor[keyKinds[k]]
			sprite.W = w
			sprite.H = h
			sprite.X = x
			sprite.Y = y
			sprite.Fix()
			s.Spotlights[k] = sprite
		}
	}
	s.combos = game.LoadNumbers(game.NumberCombo)
	s.scores = game.LoadNumbers(game.NumberScore)

	for i := range s.judgeSprite {
		src := Skin.Judge[i]
		a := game.NewAnimation([]*ebiten.Image{src})
		a.H = int(Settings.JudgeHeight * game.DisplayScale())
		scale := float64(a.H) / float64(src.Bounds().Dy())
		a.W = int(float64(src.Bounds().Dx()) * scale)
		a.X = center - a.W/2
		a.Y = int(Settings.JudgePosition*game.DisplayScale()) - a.H/2
		// a.CompositeMode = ebiten.CompositeModeSourceOver
		s.judgeSprite[i] = a
	}
	s.noteWidths = noteWidths // temp

	s.Lighting = make([]game.Animation, keyCount)
	s.LightingLN = make([]game.Animation, keyCount)
	centerXs := make([]int, keyCount)
	for k := range centerXs {
		x := center - wMiddle/2
		for k2 := 0; k2 < k; k2++ {
			x += noteWidths[k2]
		}
		x += noteWidths[k] / 2
		centerXs[k] = x
	}
	{ // suppose all frame has same size
		a := game.NewAnimation(Skin.Lighting)
		a.W = int(float64(Skin.Lighting[0].Bounds().Dx()) * Settings.LightingScale)
		a.H = int(float64(Skin.Lighting[0].Bounds().Dy()) * Settings.LightingScale)
		a.Y = int(Settings.HitPosition*game.DisplayScale()) - a.H/2
		a.CompositeMode = ebiten.CompositeModeLighter
		for k := 0; k < keyCount; k++ {
			s.Lighting[k] = a
			s.Lighting[k].X = centerXs[k] - a.W/2
		}
	}
	{
		a := game.NewAnimation(Skin.LightingLN)
		a.W = int(float64(Skin.LightingLN[0].Bounds().Dx()) * Settings.LightingLNScale)
		a.H = int(float64(Skin.LightingLN[0].Bounds().Dy()) * Settings.LightingLNScale)
		a.Y = int(Settings.HitPosition*game.DisplayScale()) - a.H/2
		a.CompositeMode = ebiten.CompositeModeLighter
		for k := 0; k < keyCount; k++ {
			s.LightingLN[k] = a
			s.LightingLN[k].X = centerXs[k] - a.W/2
		}
	}
	return *s
}

func (s *Scene) setNoteSprites() {
	keyKinds := keyKindsMap[s.chart.KeyCount]

	var wMiddle int
	for k := 0; k < s.chart.KeyCount; k++ {
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

		scale := float64(game.Settings.ScreenSize.Y) / 100
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
