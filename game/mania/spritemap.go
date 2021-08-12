package mania

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
)

// general setting 도 필요하고, mania 전용 setting도 필요하고
// general skin 도 필요하고, mania 전용 skin도 필요하고
// -> 둘다 각 settings, skin에 포함

type SpriteMapTemplate struct {
	Combo      [10]game.Sprite // unscaled
	HitResults [5]game.Sprite  // unscaled
	Stages     map[int]Stage   // 키별로 option 다름
}
type Stage struct {
	Keys    int           // todo: int8
	Notes   []game.Sprite // key
	LNHeads []game.Sprite
	LNTails []game.Sprite
	LNBodys [][]game.Sprite // key // animation
	// KeyButtons        []Sprite
	// KeyPressedButtons []Sprite
	// NoteLightings     []Sprite // unscaled
	// LNLightings       []Sprite // unscaled

	// HPBarColor Sprite // 폭맞춤x, screenHeigth
	Fixed game.Sprite
}

var SpriteMap SpriteMapTemplate

func LoadSpriteMap(skinPath string, p image.Point) {
	if !game.SpriteMap.Loaded() {
		game.LoadSpriteMap(skinPath)
	}
	loadSkin(skinPath)
	if SpriteMap.Stages == nil {
		SpriteMap.Stages = make(map[int]Stage)
	}
	// for key := range keyKinds {
	for _, key := range []int{4, 7} {
		stage := Stage{Keys: key}
		stage.Draw(p)
		SpriteMap.Stages[key] = stage
	}
}

func (s *Stage) Draw(p image.Point) { // p: screen Size
	// HPBarFrame: (폭맞춤, screenHeigth)
	scale := float64(p.Y / 100)
	noteKinds := keyKinds[s.Keys]
	noteWidths := make([]int, s.Keys&ScratchMask)
	h := int(Settings.NoteHeigth * scale)
	var fieldWidth int
	for key, kind := range noteKinds {
		w := int(Settings.NoteWidths[s.Keys&ScratchMask][kind] * scale)
		noteWidths[key] = w
		fieldWidth += w
	}
	stageOffset := StageCenter(p) - (fieldWidth)/2 // int - int. 전자는 Position일 뿐.
	{
		fixed := image.NewRGBA(image.Rectangle{image.Pt(0, 0), p})
		black := image.NewUniform(color.RGBA{0, 0, 0, 128})
		mainRect := image.Rect(stageOffset, 0, stageOffset+fieldWidth, p.Y)
		draw.Draw(fixed, mainRect, black, mainRect.Min, draw.Over)

		const hintHeight float64 = 2
		red := image.NewUniform(color.RGBA{254, 53, 53, 128})
		hintRect := image.Rect(stageOffset, int((Settings.HitPosition-hintHeight/2)*scale),
			stageOffset+fieldWidth, int((Settings.HitPosition+hintHeight/2)*scale))
		draw.Draw(fixed, hintRect, red, hintRect.Min, draw.Over)

		i, _ := ebiten.NewImageFromImage(fixed, ebiten.FilterDefault)
		s.Fixed.SetImage(i)
		s.Fixed.SetPosition(image.Point{}) // image.Pt(0, 0)
	}
	{
		s.Notes = make([]game.Sprite, len(noteWidths))
		s.LNHeads = make([]game.Sprite, len(noteWidths))
		s.LNTails = make([]game.Sprite, len(noteWidths))
		s.LNBodys = make([][]game.Sprite, len(noteWidths))
		x := stageOffset
		y := int(Settings.HitPosition*scale - float64(h)/2) // default
		for key, w := range noteWidths {
			p := image.Pt(x, y)
			x += w // for next point
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
				if Settings.LNHeadCustom {
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
				switch Settings.LNTailMode {
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
				s.LNBodys[key] = make([]game.Sprite, len(srcs))
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
		}
	}
}
