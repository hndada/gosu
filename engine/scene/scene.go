package scene

import (
	"github.com/hajimehoshi/ebiten"
)

type Scene interface {
	Ready() bool
	Update() error
	Draw(screen *ebiten.Image) // Draws scene to screen
	Close(args *Args) bool     // 모든 passed parameter는 Passed by Value.
}

// todo: ScreenSize 는 engine 에 global하게
// common unexported fields는 별 수 없이 매번 넣어줘야함

type Args struct {
	Next string // Next scene name. A scene's name can be retrieved by .(type)
	Args interface{}
}
