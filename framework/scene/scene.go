package scene

import "github.com/hndada/gosu/framework/draws"

type Scene interface {
	Update() any
	Draw(screen draws.Image)
}
