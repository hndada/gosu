package scene

import "github.com/hndada/gosu/draws"

type Scene interface {
	Update() any
	Draw(screen draws.Image)
}
