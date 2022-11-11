package draws

import "github.com/hajimehoshi/ebiten/v2"

type Box struct {
	Rectangle
	Outer Subject
	Inner [][]Box
}

func NewBox() Box {

}

func (b Box) Draw(screen *ebiten.Image) {

}
