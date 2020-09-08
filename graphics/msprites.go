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
	LNBodys [][]Sprite // key // animation
	// KeyButtons        []Sprite
	// KeyPressedButtons []Sprite
	// NoteLightings     []Sprite // unscaled
	// LNLightings       []Sprite // unscaled

	// HPBarColor Sprite // 폭맞춤x, screenHeigth
	Fixed Sprite
}

// 결국 여기서 GeoM.Scale잡아줘야됨
func (s *ManiaStage) Render(set *settings.Settings, skin skin) {
	// HPBarFrame: (폭맞춤, screenHeigth)
	scale := set.ScaleY()
	noteKinds := maniaNoteKinds[s.Keys]
	// noteSizes := make([]image.Point, s.Keys&ScratchMask)
	noteWidths := make([]int, s.Keys&ScratchMask)
	h := int(set.NoteHeigth * scale)
	var fieldWidth int
	for key, kind := range noteKinds {
		w := int(set.NoteWidths[s.Keys&ScratchMask][kind] * scale)
		noteWidths[key] = w
		fieldWidth += w
	}
	stageOffset := set.ManiaStageCenter() - (fieldWidth)/2 // int - int. 전자는 Position일 뿐.
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

		s.Fixed.i, _ = ebiten.NewImageFromImage(fixed, ebiten.FilterDefault)
		s.Fixed.p = image.Pt(0, 0)
		// s.Fixed.x, s.Fixed.y = 0, 0
		// s.Fixed.w, s.Fixed.h = set.ScreenSize().X, set.ScreenSize().Y
	}
	{
		s.Notes = make([]Sprite, len(noteWidths))
		s.LNHeads = make([]Sprite, len(noteWidths))
		s.LNTails = make([]Sprite, len(noteWidths))
		s.LNBodys = make([][]Sprite, len(noteWidths))
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
				i := skin.mania.note[noteKinds[key]]
				rw, rh := i.Size()
				op.GeoM.Scale(float64(w)/float64(rw), float64(h)/float64(rh))
				s.Notes[key].i, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
				s.Notes[key].i.DrawImage(i, op)
				s.Notes[key].p = p
			}
			{
				op := &ebiten.DrawImageOptions{}
				var i *ebiten.Image
				if set.LNHeadCustom {
					i = skin.mania.lnHead[noteKinds[key]]
				} else {
					i = skin.mania.note[noteKinds[key]]
				}
				rw, rh := i.Size()
				op.GeoM.Scale(float64(w)/float64(rw), float64(h)/float64(rh))
				s.LNHeads[key].i, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
				s.LNHeads[key].i.DrawImage(i, op)
				s.LNHeads[key].p = p
			}
			{
				op := &ebiten.DrawImageOptions{}
				var i *ebiten.Image
				switch set.LNTailMode {
				case settings.LNTailModeHead:
					i = skin.mania.lnHead[noteKinds[key]]
				case settings.LNTailModeBody:
					i, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault) // todo: test yet
				case settings.LNTailModeCustom:
					i = skin.mania.lnTail[noteKinds[key]]
				default:
					i = skin.mania.lnHead[noteKinds[key]]
				}
				rw, rh := i.Size()
				op.GeoM.Scale(float64(w)/float64(rw), float64(h)/float64(rh))
				s.LNTails[key].i, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
				s.LNTails[key].i.DrawImage(i, op)
				s.LNTails[key].p = p
			}
			{
				imgs := skin.mania.lnBody[noteKinds[key]]
				s.LNBodys[key] = make([]Sprite, len(imgs))
				for idx, i := range imgs {
					op := &ebiten.DrawImageOptions{}
					rw, rh := i.Size()
					op.GeoM.Scale(float64(w)/float64(rw), float64(h)/float64(rh))
					s.LNBodys[key][idx].i, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
					s.LNBodys[key][idx].i.DrawImage(i, op)
					s.LNBodys[key][idx].p = p
				}
			}
			// var lbsp ExpSprite
			// lbsp.vertical = true
			// lbsp.wh = sp.w
			// lbsp.x, lbsp.y = sp.x, sp.y
			// lbsp.i = skin.mania.lnBody[noteKinds[key]][0]
			// s.LNBodys[key] = make([]ExpSprite, 1)
			// s.LNBodys[key][0] = lbsp
			x += w
		}
	}
}
