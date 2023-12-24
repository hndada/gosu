package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
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

type KeyButtonOpts struct {
	ws []float64
	xs []float64
}

func NewKeyButtonOpts(keys KeysOpts) KeyButtonOpts {
	return KeyButtonOpts{
		ws: keys.ws,
		xs: keys.xs,
	}
}

func 