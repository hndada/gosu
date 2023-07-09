package scene

import (
	"github.com/hndada/gosu/draws"
)

type Scene interface {
	Update() any
	Draw(screen draws.Image)
}

type Return struct {
	From string
	To   string
	Args any
}
