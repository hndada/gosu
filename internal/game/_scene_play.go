package game

import "github.com/hajimehoshi/ebiten"


type PlayField interface {
	Update() // 다시 Game{}을 넣어야 하는데, 이는 곧 cycle을 의미.
	Draw(field *ebiten.Image)
}