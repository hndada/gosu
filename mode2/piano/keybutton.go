package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type KeyButtonRes struct {
	Imgs [2]draws.Image
}

func (kr *KeyButtonRes) Load(fsys fs.FS) {
	for i, name := range []string{"up", "down"} {
		fname := fmt.Sprintf("piano/key/%s.png", name)
		kr.Imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
}

type KeyButtonOpts struct {
}
