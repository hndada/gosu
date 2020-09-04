package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/settings"
	"image"
	"image/color"
	"image/draw"
)

// These values are applied at keys
// Example: 40 = 32 + 8 = Left-scratching 8 Key
const (
	ScratchLeft  = 1 << 5 // 32
	ScratchRight = 1 << 6 // 64
)
const ScratchMask = ^(ScratchLeft | ScratchRight)

type maniaNoteKind uint8

const (
	one maniaNoteKind = iota
	two
	middle
	pinky
)

var maniaNoteKinds = make(map[int][]maniaNoteKind)

func init() {
	maniaNoteKinds[0] = []maniaNoteKind{}
	maniaNoteKinds[1] = []maniaNoteKind{middle}
	maniaNoteKinds[2] = []maniaNoteKind{one, one}
	maniaNoteKinds[3] = []maniaNoteKind{one, middle, one}
	maniaNoteKinds[4] = []maniaNoteKind{one, two, two, one}
	maniaNoteKinds[5] = []maniaNoteKind{one, two, middle, two, one}
	maniaNoteKinds[6] = []maniaNoteKind{one, two, one, one, two, one}
	maniaNoteKinds[7] = []maniaNoteKind{one, two, one, middle, one, two, one}
	maniaNoteKinds[8] = []maniaNoteKind{pinky, one, two, one, one, two, one, pinky}
	maniaNoteKinds[9] = []maniaNoteKind{pinky, one, two, one, middle, one, two, one, pinky}
	maniaNoteKinds[10] = []maniaNoteKind{pinky, one, two, one, middle, middle, one, two, one, pinky}

	for i := 1; i <= 8; i++ { // 정말 잘 짠듯
		maniaNoteKinds[i|ScratchLeft] = append([]maniaNoteKind{pinky}, maniaNoteKinds[i-1]...)
		maniaNoteKinds[i|ScratchRight] = append(maniaNoteKinds[i-1], pinky)
	}
}

type ManiaStage struct {
	Keys    int      // todo: int8
	Notes   []Sprite // key
	LNHeads []Sprite
	LNTails []Sprite
	LNBodys [][]ExpSprite // key // animation
	// KeyButtons        []Sprite
	// KeyPressedButtons []Sprite
	// NoteLightings     []Sprite // unscaled
	// LNLightings       []Sprite // unscaled

	// HPBarColor Sprite // 폭맞춤x, screenHeigth
	Fixed Sprite
}

func (s *ManiaStage) Render(set *settings.Settings, skin skin) {
	scale := set.ScaleY()
	noteKinds := maniaNoteKinds[s.Keys]
	noteSizes := make([]image.Point, s.Keys&ScratchMask)
	{
		h := int(set.NoteHeigth * scale)
		for key, kind := range noteKinds {
			w := int(set.NoteWidths[s.Keys&ScratchMask][kind] * scale)
			noteSizes[key] = image.Pt(w, h)
		}
	}

	// HPBarFrame: (폭맞춤, screenHeigth)
	var fieldWidth int
	for _, ns := range noteSizes {
		fieldWidth += ns.X
	}
	stageOffset := set.ManiaStageCenter() - (fieldWidth)/2 // int - int. 전자는 Position일 뿐.

	fixed := image.NewRGBA(image.Rect(0, 0, set.ScreenSize().X, set.ScreenSize().Y))
	black := image.NewUniform(color.RGBA{0, 0, 0, 128})
	mainRect := image.Rect(stageOffset, 0, stageOffset+fieldWidth, set.ScreenSize().Y)
	draw.Draw(fixed, mainRect, black, mainRect.Min, draw.Over)

	red := image.NewUniform(color.RGBA{254, 106, 109, 128})
	hintRect := image.Rect(stageOffset, int(set.HitPosition*scale),
		stageOffset+fieldWidth, int(set.HitPosition*scale+set.NoteHeigth*scale))
	draw.Draw(fixed, hintRect, red, hintRect.Min, draw.Over)

	s.Fixed.i, _ = ebiten.NewImageFromImage(fixed, ebiten.FilterDefault)
	s.Fixed.x, s.Fixed.y = 0, 0
	s.Fixed.w, s.Fixed.h = set.ScreenSize().X, set.ScreenSize().Y

	s.Notes = make([]Sprite, len(noteSizes))
	s.LNHeads = make([]Sprite, len(noteSizes))
	s.LNTails = make([]Sprite, len(noteSizes))
	s.LNBodys = make([][]ExpSprite, len(noteSizes))
	{
		hitPosition := int(set.HitPosition * set.ScaleY())
		var x int
		for key, ns := range noteSizes {
			var sp Sprite
			sp.w, sp.h = ns.X, ns.Y
			sp.x = x + stageOffset
			sp.y = hitPosition
			x += sp.w

			nsp := sp
			nsp.i = skin.mania.note[noteKinds[key]]
			s.Notes[key] = nsp

			lhsp := sp
			if set.LNHeadCustom {
				lhsp.i = skin.mania.lnHead[noteKinds[key]]
			} else {
				lhsp.i = skin.mania.note[noteKinds[key]]
			}
			s.LNHeads[key] = lhsp

			ltsp := sp
			switch set.LNTailMode {
			case settings.LNTailModeHead:
				ltsp.i = skin.mania.lnHead[noteKinds[key]]
			case settings.LNTailModeBody:
				ltsp.i, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault) // todo: test yet
			case settings.LNTailModeCustom:
				ltsp.i = skin.mania.lnTail[noteKinds[key]]
			default:
				ltsp.i = skin.mania.lnHead[noteKinds[key]]
			}
			s.LNTails[key] = ltsp

			var lbsp ExpSprite
			lbsp.vertical = true
			lbsp.wh = sp.w
			lbsp.x, lbsp.y = sp.x, sp.y
			lbsp.i = skin.mania.lnBody[noteKinds[key]][0]
			s.LNBodys[key] = make([]ExpSprite, 1)
			s.LNBodys[key][0] = lbsp
		}
	}
}
