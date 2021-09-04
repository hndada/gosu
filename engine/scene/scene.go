package scene

import (
	"github.com/hajimehoshi/ebiten"
)

// common unexported fields는 별 수 없이 매번 넣어줘야함
type Scene interface {
	Ready() bool
	Update() error
	Draw(screen *ebiten.Image) // Draws scene to screen
	Close(args *Args) bool     // 모든 passed parameter는 Passed by Value.
}

type Args struct {
	Next string // Next scene name. A scene's name can be retrieved by .(type)
	Args interface{}
}
