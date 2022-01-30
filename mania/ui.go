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
	sUI := new(sceneUI)
	scale := float64(common.Settings.ScreenSizeY) / 100
	keyKinds := keyKindsMap[WithScratch(keyCount)]
	unscaledNoteWidths := Settings.NoteWidths[keyCount]

	noteWidths := make([]int, keyCount)
	for key, kind := range keyKinds {
		noteWidths[key] = int(unscaledNoteWidths[kind] * scale)
	}
	playfieldImage := ebiten.NewImage(common.Settings.ScreenSizeX, common.Settings.ScreenSizeY)

	p := Settings.StagePosition / 100 // proportion
	center := int(float64(common.Settings.ScreenSizeX) * p)
	var wLeft, wMiddle int
	{ // main
		for _, nw := range noteWidths {
			wMiddle += nw
		}
		h := common.Settings.ScreenSizeY

		// seems ebiten's Fill() doesn't accept alpha value
		mainSrc := image.NewRGBA(image.Rect(0, 0, wMiddle, h))
		r := image.Rectangle{image.ZP, playfieldImage.Bounds().Size()}
		draw.Draw(mainSrc, r, &image.Uniform{black}, image.ZP, draw.Over)
		main := ebiten.NewImageFromImage(mainSrc)

		x := center - wMiddle/2 // int - int
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		op.ColorM.Scale(0, 0, 0, 1)
		op.ColorM.ChangeHSV(0, 1, Settings.PlayfieldDimness)
		playfieldImage.DrawImage(main, op)
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
		playfieldImage.DrawImage(hint, op)
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
		h := common.Settings.ScreenSizeY
		scale := float64(h) / float64(src.Bounds().Dy())
		wLeft = int(float64(src.Bounds().Dx()) * scale)
		x := center - wMiddle/2 - wLeft
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
		playfieldImage.DrawImage(src, op)
	}
	{
		src := Skin.StageRight
		h := common.Settings.ScreenSizeY
		scale := float64(h) / float64(src.Bounds().Dy())
		// wRight = int(float64(src.Bounds().Dx()) * scale)
		x := center + wMiddle/2
		y := 0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
		playfieldImage.DrawImage(src, op)
	}
	{ // Beware of setting WHXY: the image goes 90-degree rotating
		src := Skin.HPBar
		sp := ui.NewSprite(src)
		sp.H = int(Settings.HPHeight * common.DisplayScale())
		scale := float64(sp.H) / float64(src.Bounds().Dy())
		sp.W = int(float64(src.Bounds().Dx()) * scale)
		sp.X = center + wMiddle/2
		sp.Y = common.Settings.ScreenSizeY - sp.H
		sUI.HPBar = ui.NewFixedSprite(sp)
	}
	{ // Its size can be different with HP Bar Image's.
		src := Skin.HPBarColor
		sp := ui.NewSprite(src)
		sp.H = int(Settings.HPHeight * common.DisplayScale())
		scale := float64(sp.H) / float64(src.Bounds().Dy())
		sp.W = int(float64(src.Bounds().Dx()) * scale)
		sp.X = center + wMiddle/2 // + s.HPBar.W/2
		sp.Y = common.Settings.ScreenSizeY - sp.H
		// y := int(Settings.HitPosition*common.DisplayScale()) - h
		sUI.HPBarColor = ui.NewFixedSprite(sp)

		sp2 := ui.NewSprite(ebiten.NewImage(sp.W, sp.H))
		sp2.W = sp.W
		sp2.H = 0 // HP:100
		sp2.X = sp.X
		sp2.Y = sp.Y
		sp2.CompositeMode = ebiten.CompositeModeSourceOut
		sUI.HPBarMask = sp2
	}
	playfieldSprite := ui.NewSprite(playfieldImage)
	playfieldSprite.W = common.Settings.ScreenSizeX
	playfieldSprite.H = common.Settings.ScreenSizeY
	playfieldSprite.X = 0
	playfieldSprite.Y = 0
	sUI.playfield = ui.NewFixedSprite(playfieldSprite)

	sUI.stageKeys = make([]ui.FixedSprite, keyCount)
	sUI.stageKeysPressed = make([]ui.FixedSprite, keyCount)

	// Each skin has own empty space.
	for k := 0; k < keyCount; k++ {
		sp := ui.NewSprite(Skin.StageKeys[keyKinds[k]])
		sp.W = noteWidths[k]      // Note widths can be different, while its source image size is same.
		sp.X = center - wMiddle/2 // int - int
		for k2 := 0; k2 < k; k2++ {
			sp.X += noteWidths[k2]
		}
		sp.Y = int(Settings.HitPosition * common.DisplayScale()) // + hHint/2
		sp.H = common.Settings.ScreenSizeY - sp.Y
		sUI.stageKeys[k] = ui.NewFixedSprite(sp)

		sp2 := sp
		src2 := Skin.StageKeysPressed[keyKinds[k]]
		sp2.SetImage(src2)
		sUI.stageKeysPressed[k] = ui.NewFixedSprite(sp2)
	}
	{
		src := Skin.StageLight
		sp := ui.NewSprite(src)
		sUI.Spotlights = make([]ui.FixedSprite, keyCount)
		for k := 0; k < keyCount; k++ {
			w := noteWidths[k] // Note widths can be different, while its source image size is same.
			scale := float64(w) / float64(src.Bounds().Size().X)
			sp.H = int(float64(src.Bounds().Size().Y) * scale)
			sp.X = center - wMiddle/2 // int - int
			for k2 := 0; k2 < k; k2++ {
				sp.X += noteWidths[k2]
			}
			sp.Y = int(Settings.HitPosition*common.DisplayScale()) - sp.H
			sp.Color = Settings.SpotlightColor[keyKinds[k]]
			sp.W = w
			sUI.Spotlights[k] = ui.NewFixedSprite(sp)
		}
	}
	sUI.combos = common.LoadNumbers(common.NumberCombo)
	sUI.scores = common.LoadNumbers(common.NumberScore)

	for i := range sUI.judgeSprite {
		src := Skin.Judge[i]
		a := ui.NewAnimation([]*ebiten.Image{src})
		a.H = int(Settings.JudgeHeight * common.DisplayScale())
		scale := float64(a.H) / float64(src.Bounds().Dy())
		a.W = int(float64(src.Bounds().Dx()) * scale)
		a.X = center - a.W/2
		a.Y = int(Settings.JudgePosition*common.DisplayScale()) - a.H/2
		// a.CompositeMode = ebiten.CompositeModeSourceOver
		sUI.judgeSprite[i] = a
	}
	sUI.noteWidths = noteWidths // temp

	sUI.Lighting = make([]ui.Animation, keyCount)
	sUI.LightingLN = make([]ui.Animation, keyCount)
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
			sUI.Lighting[k] = a
			sUI.Lighting[k].X = centerXs[k] - a.W/2
		}
	}
	{
		a := ui.NewAnimation(Skin.LightingLN)
		a.W = int(float64(Skin.LightingLN[0].Bounds().Dx()) * Settings.LightingLNScale)
		a.H = int(float64(Skin.LightingLN[0].Bounds().Dy()) * Settings.LightingLNScale)
		a.Y = int(Settings.HitPosition*common.DisplayScale()) - a.H/2
		a.CompositeMode = ebiten.CompositeModeLighter
		for k := 0; k < keyCount; k++ {
			sUI.LightingLN[k] = a
			sUI.LightingLN[k].X = centerXs[k] - a.W/2
		}
	}
	return *sUI
}

func (s *Scene) setNoteSprites() {
	keyKinds := keyKindsMap[WithScratch(s.chart.KeyCount)]

	var wMiddle int
	for k := 0; k < s.chart.KeyCount; k++ {
		wMiddle += s.noteWidths[k]
	}
	xStart := (common.Settings.ScreenSizeX - wMiddle) / 2
	for i, n := range s.chart.Notes {
		var sprite ui.Sprite
		kind := keyKinds[n.key]
		switch n.Type {
		case TypeNote, TypeLNTail: // temp
			sprite = ui.NewSprite(Skin.Note[kind])
		case TypeLNHead:
			sprite = ui.NewSprite(Skin.LNHead[kind])
		}

		scale := float64(common.Settings.ScreenSizeY) / 100
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
