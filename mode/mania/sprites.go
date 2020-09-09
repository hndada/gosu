package mania

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mode"
	"image"
	"image/color"
	"image/draw"
)

// general setting 도 필요하고, mania 전용 setting도 필요하고
// general skin 도 필요하고, mania 전용 skin도 필요하고
// -> 둘다 각 settings, skin에 포함
type Sprites struct {
	settings   *Settings
	skin       *skin
	Combo      [10]mode.Sprite // unscaled
	HitResults [5]mode.Sprite  // unscaled
	Stages     map[int]Stage   // 키별로 option 다름
}

type Stage struct {
	Keys    int           // todo: int8
	Notes   []mode.Sprite // key
	LNHeads []mode.Sprite
	LNTails []mode.Sprite
	LNBodys [][]mode.Sprite // key // animation
	// KeyButtons        []Sprite
	// KeyPressedButtons []Sprite
	// NoteLightings     []Sprite // unscaled
	// LNLightings       []Sprite // unscaled

	// HPBarColor Sprite // 폭맞춤x, screenHeigth
	Fixed mode.Sprite
}

func (s *Sprites) Render(settings *Settings) {
	s.settings = settings
	// todo: 깔끔하게
	if s.skin == nil {
		s.skin = &skin{}
		s.skin.load(`C:\Users\hndada\Documents\GitHub\hndada\gosu\test\Skin`)
	}
	if s.Stages == nil {
		s.Stages = make(map[int]Stage)
	}
	// for key := range keyKinds {
	for _, key := range []int{4, 7} {
		stage := Stage{Keys: key}
		stage.Render(s.settings, s.skin)
		s.Stages[key] = stage
	}
}

func (s *Stage) Render(set *Settings, skin *skin) {
	// HPBarFrame: (폭맞춤, screenHeigth)
	scale := set.ScaleY()
	noteKinds := keyKinds[s.Keys]
	// noteSizes := make([]image.Point, s.Keys&ScratchMask)
	noteWidths := make([]int, s.Keys&ScratchMask)
	h := int(set.NoteHeigth * scale)
	var fieldWidth int
	for key, kind := range noteKinds {
		w := int(set.NoteWidths[s.Keys&ScratchMask][kind] * scale)
		noteWidths[key] = w
		fieldWidth += w
	}
	stageOffset := set.StageCenter(set.ScreenSize()) - (fieldWidth)/2 // int - int. 전자는 Position일 뿐.
	{
		fixed := image.NewRGBA(image.Rectangle{image.Pt(0, 0), set.ScreenSize()})
		black := image.NewUniform(color.RGBA{0, 0, 0, 128})
		mainRect := image.Rect(stageOffset, 0, stageOffset+fieldWidth, set.ScreenSize().Y)
		draw.Draw(fixed, mainRect, black, mainRect.Min, draw.Over)

		const hintHeight float64 = 2
		red := image.NewUniform(color.RGBA{254, 53, 53, 128})
		hintRect := image.Rect(stageOffset, int((set.HitPosition-hintHeight/2)*scale),
			stageOffset+fieldWidth, int((set.HitPosition+hintHeight/2)*scale))
		draw.Draw(fixed, hintRect, red, hintRect.Min, draw.Over)

		i, _ := ebiten.NewImageFromImage(fixed, ebiten.FilterDefault)
		s.Fixed.SetImage(i)
		s.Fixed.SetPosition(image.Pt(0, 0))
		// s.Fixed.x, s.Fixed.y = 0, 0
		// s.Fixed.w, s.Fixed.h = set.ScreenSize().X, set.ScreenSize().Y
	}
	{
		s.Notes = make([]mode.Sprite, len(noteWidths))
		s.LNHeads = make([]mode.Sprite, len(noteWidths))
		s.LNTails = make([]mode.Sprite, len(noteWidths))
		s.LNBodys = make([][]mode.Sprite, len(noteWidths))
		x := stageOffset
		y := int(set.HitPosition*scale - float64(h)/2) // default
		for key, w := range noteWidths {
			// var sp Sprite
			// sp.w, sp.h = size.X, size.Y // 이미 w, h 맞춰서 나옴
			// sp.x = x + stageOffset
			// sp.y = hitPosition
			p := image.Pt(x, y)
			{
				op := &ebiten.DrawImageOptions{}
				src := skin.note[noteKinds[key]]
				rw, rh := src.Size()
				op.GeoM.Scale(float64(w)/float64(rw), float64(h)/float64(rh))
				dst, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
				dst.DrawImage(src, op)
				s.Notes[key].SetImage(dst)
				s.Notes[key].SetPosition(p)
			}
			{
				op := &ebiten.DrawImageOptions{}
				var src *ebiten.Image
				if set.LNHeadCustom {
					src = skin.lnHead[noteKinds[key]]
				} else {
					src = skin.note[noteKinds[key]]
				}
				rw, rh := src.Size()
				op.GeoM.Scale(float64(w)/float64(rw), float64(h)/float64(rh))
				dst, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
				dst.DrawImage(src, op)
				s.LNHeads[key].SetImage(dst)
				s.LNHeads[key].SetPosition(p)
			}
			{
				op := &ebiten.DrawImageOptions{}
				var src *ebiten.Image
				switch set.LNTailMode {
				case LNTailModeHead:
					src = skin.lnHead[noteKinds[key]]
				case LNTailModeBody:
					src, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault) // todo: test yet
				case LNTailModeCustom:
					src = skin.lnTail[noteKinds[key]]
				default:
					src = skin.lnHead[noteKinds[key]]
				}
				rw, rh := src.Size()
				op.GeoM.Scale(float64(w)/float64(rw), float64(h)/float64(rh))
				dst, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
				dst.DrawImage(src, op)
				s.LNTails[key].SetImage(dst)
				s.LNTails[key].SetPosition(p)
			}
			{
				srcs := skin.lnBody[noteKinds[key]]
				s.LNBodys[key] = make([]mode.Sprite, len(srcs))
				for idx, src := range srcs {
					op := &ebiten.DrawImageOptions{}
					rw, rh := src.Size()
					ratio := float64(w) / float64(rw)
					op.GeoM.Scale(ratio, ratio)

					rh2 := int(float64(rh) * ratio)
					count := 4000 / rh2
					bufferedHeight := rh2 * count
					dst, _ := ebiten.NewImage(w, bufferedHeight, ebiten.FilterDefault)
					for c := 0; c < count; c++ {
						dst.DrawImage(src, op)
						op.GeoM.Translate(0, float64(rh2))
					}
					s.LNBodys[key][idx].SetImage(dst)
					s.LNBodys[key][idx].SetPosition(p)
				}
			}
			// var lbsp ExpSprite
			// lbsp.vertical = true
			// lbsp.wh = sp.w
			// lbsp.x, lbsp.y = sp.x, sp.y
			// lbsp.i = skin.lnBody[noteKinds[key]][0]
			// s.LNBodys[key] = make([]ExpSprite, 1)
			// s.LNBodys[key][0] = lbsp
			x += w
		}
	}
}
