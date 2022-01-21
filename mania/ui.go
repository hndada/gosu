package mania

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/ui"
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
	noteWidths       []int // TODO: It is currently used only at setNoteSprites()
	playfield        ui.FixedSprite
	stageKeys        []ui.FixedSprite
	stageKeysPressed []ui.FixedSprite

	combos      [10]ui.Sprite
	scores      [10]ui.Sprite
	judgeSprite [len(Judgments)]ui.Animation // TODO: rename
	Spotlights  []ui.FixedSprite             // Blinking component when pressing keys

	HPBar      ui.FixedSprite // it can be in playfield
	HPBarColor ui.FixedSprite // actually, it can also go to playfield
	HPBarMask  ui.Sprite
	hpScreen   *ebiten.Image

	Lighting   []ui.Animation // It should be able to be drawn simultaneously in all lanes
	LightingLN []ui.Animation
}

// A width of screen size doesn't affect to UI size; only height does: standard is 100
func newSceneUI(keyCount int) sceneUI {
	s := new(sceneUI)
	scale := float64(common.Settings.ScreenSize.Y) / 100
	keyKinds := keyKindsMap[WithScratch(keyCount)]
	unscaledNoteWidths := Settings.NoteWidths[keyCount]

	noteWidths := make([]int, keyCount)
	for key, kind := range keyKinds {
		noteWidths[key] = int(unscaledNoteWidths[kind] * scale)
	}
	i := ebiten.NewImage(common.Settings.ScreenSize.X, common.Settings.ScreenSize.Y)

	p := Settings.StagePosition / 100 // proportion
	center := int(float64(common.Settings.ScreenSize.X) * p)
	var wLeft, wMiddle int
	{ // main
		for _, nw := range noteWidths {
			wMiddle += nw
		}
		h := common.Settings.ScreenSize.Y

		// seems ebiten's Fill() doesn't accept alpha value
		mainSrc := image.NewRGBA(image.Rect(0, 0, wMiddle, h))
		r := image.Rectangle{image.ZP, i.Bounds().Size()}
		draw.Draw(mainSrc, r, &image.Uniform{black}, image.ZP, draw.Over)
		main := ebiten.NewImageFromImage(mainSrc)

		x := center - wMiddle/2 // int - int
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.ColorM.Scale(0, 0, 0, 1)
		op.ColorM.ChangeHSV(0, 1, Settings.PlayfieldDimness)
		i.DrawImage(main, op)
	}
	// Important: There's no guarantee that judge-line locates at the very bottom at 'mania-stage-hint' image.
	// cf. 'mania-stage-bottom'

	// var hHint int
	{ // no-skin ver

		h := int(Settings.JudgeLineHeight * common.DisplayScale())
		hint := ebiten.NewImage(wMiddle, h)
		hint.Fill(red)

		x := center - wMiddle/2 // int - int
		y := int(Settings.HitPosition*common.DisplayScale()) - h
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		i.DrawImage(hint, op)
	}
	// {
	// 	src := Skin.StageHint
	// 	scale := float64(wMiddle) / float64(src.Bounds().Dx())
	// 	h := int(float64(src.Bounds().Dy()) * scale)
	// 	x := center - wMiddle/2
	// 	y := int(Settings.HitPosition*common.DisplayScale()) - h
	// 	op := &ebiten.DrawImageOptions{}
	// 	op.GeoM.Scale(scale, scale)
	// 	op.GeoM.Translate(float64(x), float64(y))
	// 	i.DrawImage(src, op)
	// 	// hHint = h
	// }
	{
		src := Skin.StageLeft
		h := common.Settings.ScreenSize.Y
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
		h := common.Settings.ScreenSize.Y
		scale := float64(h) / float64(src.Bounds().Dy())
		// wRight = int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
		i.DrawImage(src, op)
	}
	{ // Beware of setting WHXY: the image goes 90-degree rotating
		src := Skin.HPBar
		sprite := ui.NewFixedSprite(src)
		h := int(Settings.HPHeight * common.DisplayScale())
		scale := float64(h) / float64(src.Bounds().Dy())
		w := int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2
		y := common.Settings.ScreenSize.Y - h
		sprite.W = w
		sprite.H = h
		sprite.X = x
		sprite.Y = y
		sprite.Fix()
		s.HPBar = sprite
	}
	{ // Its size can be different with HP Bar Image's.
		src := Skin.HPBarColor
		sprite := ui.NewFixedSprite(src)
		h := int(Settings.HPHeight * common.DisplayScale())
		scale := float64(h) / float64(src.Bounds().Dy())
		w := int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2 // + s.HPBar.W/2
		y := common.Settings.ScreenSize.Y - h
		// y := int(Settings.HitPosition*common.DisplayScale()) - h
		sprite.W = w
		sprite.H = h
		sprite.X = x
		sprite.Y = y
		sprite.Fix()
		s.HPBarColor = sprite

		mask := ebiten.NewImage(w, h)
		sprite2 := ui.NewSprite(mask)
		sprite2.W = w
		sprite2.H = 0 // HP:100
		sprite2.X = x
		sprite2.Y = y
		sprite2.CompositeMode = ebiten.CompositeModeSourceOut
		s.HPBarMask = sprite2
	}
	s.playfield = ui.NewFixedSprite(i)
	s.playfield.W = common.Settings.ScreenSize.X
	s.playfield.H = common.Settings.ScreenSize.Y
	s.playfield.X = 0
	s.playfield.Y = 0
	s.playfield.Fix() // TODO: add BG at here

	s.stageKeys = make([]ui.FixedSprite, keyCount)
	s.stageKeysPressed = make([]ui.FixedSprite, keyCount)

	// Each skin has own empty space.
	for k := 0; k < keyCount; k++ {
		var sprite ui.FixedSprite
		src := Skin.StageKeys[keyKinds[k]]
		sprite = ui.NewFixedSprite(src)

		w := noteWidths[k] // Note widths can be different, while its source image size is same.

		// scale := float64(sprite.W) / float64(src.Bounds().Size().X)
		// sprite.H = int(float64(src.Bounds().Size().Y) * scale)
		x := center - wMiddle/2 // int - int
		for k2 := 0; k2 < k; k2++ {
			x += noteWidths[k2]
		}
		y := int(Settings.HitPosition * common.DisplayScale()) // + hHint/2
		// fmt.Println(hHint)
		// y := int((Settings.HitPosition - Settings.NoteHeigth/2 -
		// 	4*Settings.NoteHeigth/2) * common.DisplayScale()) // TODO: why?
		h := common.Settings.ScreenSize.Y - y

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
		sprite := ui.NewFixedSprite(src)
		s.Spotlights = make([]ui.FixedSprite, keyCount)
		for k := 0; k < keyCount; k++ {
			w := noteWidths[k] // Note widths can be different, while its source image size is same.
			scale := float64(w) / float64(src.Bounds().Size().X)
			h := int(float64(src.Bounds().Size().Y) * scale)
			x := center - wMiddle/2 // int - int
			for k2 := 0; k2 < k; k2++ {
				x += noteWidths[k2]
			}
			y := int(Settings.HitPosition*common.DisplayScale()) - h
			sprite.Color = Settings.SpotlightColor[keyKinds[k]]
			sprite.W = w
			sprite.H = h
			sprite.X = x
			sprite.Y = y
			sprite.Fix()
			s.Spotlights[k] = sprite
		}
	}
	s.combos = common.LoadNumbers(common.NumberCombo)
	s.scores = common.LoadNumbers(common.NumberScore)

	for i := range s.judgeSprite {
		src := Skin.Judge[i]
		a := ui.NewAnimation([]*ebiten.Image{src})
		a.H = int(Settings.JudgeHeight * common.DisplayScale())
		scale := float64(a.H) / float64(src.Bounds().Dy())
		a.W = int(float64(src.Bounds().Dx()) * scale)
		a.X = center - a.W/2
		a.Y = int(Settings.JudgePosition*common.DisplayScale()) - a.H/2
		// a.CompositeMode = ebiten.CompositeModeSourceOver
		s.judgeSprite[i] = a
	}
	s.noteWidths = noteWidths // temp

	s.Lighting = make([]ui.Animation, keyCount)
	s.LightingLN = make([]ui.Animation, keyCount)
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
		a := ui.NewAnimation(Skin.Lighting)
		a.W = int(float64(Skin.Lighting[0].Bounds().Dx()) * Settings.LightingScale)
		a.H = int(float64(Skin.Lighting[0].Bounds().Dy()) * Settings.LightingScale)
		a.Y = int(Settings.HitPosition*common.DisplayScale()) - a.H/2
		a.CompositeMode = ebiten.CompositeModeLighter
		for k := 0; k < keyCount; k++ {
			s.Lighting[k] = a
			s.Lighting[k].X = centerXs[k] - a.W/2
		}
	}
	{
		a := ui.NewAnimation(Skin.LightingLN)
		a.W = int(float64(Skin.LightingLN[0].Bounds().Dx()) * Settings.LightingLNScale)
		a.H = int(float64(Skin.LightingLN[0].Bounds().Dy()) * Settings.LightingLNScale)
		a.Y = int(Settings.HitPosition*common.DisplayScale()) - a.H/2
		a.CompositeMode = ebiten.CompositeModeLighter
		for k := 0; k < keyCount; k++ {
			s.LightingLN[k] = a
			s.LightingLN[k].X = centerXs[k] - a.W/2
		}
	}
	return *s
}

func (s *Scene) setNoteSprites() {
	keyKinds := keyKindsMap[WithScratch(s.chart.KeyCount)]

	var wMiddle int
	for k := 0; k < s.chart.KeyCount; k++ {
		wMiddle += s.noteWidths[k]
	}
	xStart := (common.Settings.ScreenSize.X - wMiddle) / 2
	for i, n := range s.chart.Notes {
		var sprite ui.Sprite
		kind := keyKinds[n.key]
		switch n.Type {
		case TypeNote, TypeLNTail: // temp
			sprite = ui.NewSprite(Skin.Note[kind])
		case TypeLNHead:
			sprite = ui.NewSprite(Skin.LNHead[kind])
		}

		scale := float64(common.Settings.ScreenSize.Y) / 100
		sprite.H = int(Settings.NoteHeigth * scale)
		sprite.W = s.noteWidths[n.key]
		x := xStart
		for k := 0; k < n.key; k++ {
			x += s.noteWidths[k]
		}
		sprite.X = x
		y := Settings.HitPosition - n.position*s.speed - float64(sprite.H)/2
		sprite.Y = int(y * scale)
		s.chart.Notes[i].Sprite = sprite
	}

	// LN body sprite
	// All sprites should be connected with objects which update sprites' value
	for i, tail := range s.chart.Notes {
		if tail.Type != TypeLNTail {
			continue
		}
		head := s.chart.Notes[tail.prev]
		ls := ui.LongSprite{
			Vertical: true,
		}
		ls.SetImage(Skin.LNBody[keyKinds[tail.key]]) // temp: no animation support
		ls.W = tail.Sprite.W
		ls.H = head.Sprite.Y - tail.Sprite.Y
		ls.X = tail.Sprite.X
		ls.Y = tail.Sprite.Y
		ls.Saturation = 1
		ls.Dimness = 1
		s.chart.Notes[i].LongSprite = ls
	}
}
