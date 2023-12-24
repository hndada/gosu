package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

// All names of fields in Asset ends with their types.

type KeyButtonsRes struct {
	imgs [2]draws.Image
}

func (kr *KeyButtonsRes) Load(fsys fs.FS) {
	for i, name := range []string{"up", "down"} {
		fname := fmt.Sprintf("piano/key/%s.png", name)
		kr.imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
}

type KeyButtonsOpts struct {
	ws []float64
	h  float64
	xs []float64
	y  float64 // center top
}

func NewKeyButtonsOpts(keys KeysOpts) KeyButtonsOpts {
	return KeyButtonsOpts{
		ws: keys.ws,
		h:  mode.ScreenH - keys.BaselineY,
		xs: keys.xs,
		y:  keys.BaselineY,
	}
}

type KeyButtonComp struct {
	spritesList [][2]draws.Sprite
}

func NewKeyButtonComp(res KeyButtonsRes, opts KeyButtonsOpts) (comp KeyButtonComp) {
	comp.spritesList = make([][2]draws.Sprite, len(opts.ws))
	for k := range comp.spritesList {
		for i, img := range res.imgs {
			sprite := draws.NewSprite(img)
			sprite.SetSize(opts.ws[k], opts.h)
			sprite.Locate(opts.xs[k], opts.y, draws.CenterTop)
			comp.spritesList[k][i] = sprite
		}
	}
	return
}
