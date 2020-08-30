package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/settings"
	"image"
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
	s.Notes = make([]Sprite, len(noteSizes))
	s.LNHeads = make([]Sprite, len(noteSizes))
	s.LNTails = make([]Sprite, len(noteSizes))
	{
		hitPosition := int(set.HitPosition * set.ScaleY())
		var x int
		for key, ns := range noteSizes {
			var sp Sprite
			sp.w, sp.h = ns.X, ns.Y
			sp.x = x
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
		}
	}
	// todo: fixed stage
	// 대충 그리고 test 해보기
	// main *ebiten.image // fieldWidth, screenHeight (generated)
	// HPBarFrame      Sprite     // 폭맞춤x, screenHeigth
	var fieldWidth int
	for _, ns := range noteSizes {
		fieldWidth += ns.X
	}
	stageOffset := set.ManiaStageCenter() - (fieldWidth)/2 // int - int. 전자는 Position일 뿐.
}
