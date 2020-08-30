package config

import (
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
	Keys             int      // todo: int8
	Note             []Sprite // key
	LNHead           []Sprite
	LNBody           [][]Sprite // key // animation
	LNTail           []Sprite
	KeyButton        []Sprite
	KeyButtonPressed []Sprite
	NoteLighting     []Sprite // unscaled
	LNLighting       []Sprite // unscaled

	HPBarColor Sprite // 폭맞춤x, screenHeigth
	Fixed      Sprite
}

func (s *ManiaStage) Render(settings *Settings) {
	// main *ebiten.image // fieldWidth, screenHeight (generated)
	// HPBarFrame      Sprite     // 폭맞춤x, screenHeigth
	scale := float64(settings.screenSize.Y) / 100
	noteSize := make([]image.Point, s.Keys&ScratchMask)
	{
		h := int(settings.NoteHeigth * scale)
		for key, kind := range maniaNoteKinds[s.Keys] {
			w := int(settings.NoteWidths[s.Keys][kind] * scale)
			noteSize[key] = image.Pt(w, h)
		}
	}
	s.RenderLNTail(noteSize, settings.LNTailMode)
	var fieldWidth int
	for _, ns := range noteSize {
		fieldWidth += ns.X
	}
	stageOffset := settings.maniaStageCenter() - (fieldWidth)/2 // int - int. 전자는 Position일 뿐.
}

func (s *ManiaStage) RenderLNTail(noteSize []image.Point, mode uint8) {
	switch mode {
	case LNTailModeHead:
	case LNTailModeBody:
	case LNTailModeCustom:
	default:
		// 헤드 이미지 로드
	}
}
