package draws

import "github.com/hajimehoshi/ebiten/v2"

func (s Sprite) IsMouseIn() bool {
	p := NewVector2FromInts(ebiten.CursorPosition())
	p = p.Sub(s.Min())
	max := s.Max()
	return 0 <= p.X && p.X <= max.X && 0 <= p.Y && p.Y <= max.Y
}

func (s Sprite) isMouseClick(mb ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(mb) && s.IsMouseIn()
}
func (s Sprite) IsMouseLeftClick() bool {
	if !s.IsMouseIn() {
		return false
	}
	return s.isMouseClick(ebiten.MouseButtonLeft)
}
func (s Sprite) IsMouseMiddleClick() bool {
	if !s.IsMouseIn() {
		return false
	}
	return s.isMouseClick(ebiten.MouseButtonMiddle)
}
func (s Sprite) IsMouseRightClick() bool {
	if !s.IsMouseIn() {
		return false
	}
	return s.isMouseClick(ebiten.MouseButtonRight)
}
