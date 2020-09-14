package mania

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

// todo: Scored를 여기에 넣을지 Note에 넣을지 // 판정 완료 여부. 명암 변화 및 tail에선 점수에까지 영향
type NoteSprite struct {
	h        int     // same with second value of NoteSprite.i.Size()
	x        float64 // x is fixed among mania notes
	position float64 // positive value // todo: Note로
	y        float64
	i        *ebiten.Image
	op       *ebiten.DrawImageOptions
}

type LNSprite struct {
	head   *NoteSprite
	tail   *NoteSprite
	length float64 // todo: Note로

	height float64
	i      *ebiten.Image
	bodyop *ebiten.DrawImageOptions
}

func (s LNSprite) DrawLN(screen *ebiten.Image) {
	_, h := s.i.Size()
	count, remainder := int(s.height)/h, int(s.height)%h+1
	s.bodyop.GeoM.Reset()
	s.bodyop.GeoM.Translate(s.tail.x, s.tail.y)

	firstRect := s.i.Bounds()
	firstRect.Min = image.Pt(0, h-remainder)
	screen.DrawImage(s.i.SubImage(firstRect).(*ebiten.Image), s.bodyop)
	s.bodyop.GeoM.Translate(0, float64(remainder))

	for c := 0; c < count; c++ {
		screen.DrawImage(s.i, s.bodyop)
		s.bodyop.GeoM.Translate(0, float64(h))
	}
}
