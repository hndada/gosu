package piano

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

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

// Put suffix 'List' when suffix 's' is not available.
type KeyButtonsComp struct {
	keyDowns    []bool
	spritesList [][2]draws.Sprite
	startTimes  []time.Time
	minDuration int32 // milliseconds
}

func NewKeyButtonsComp(res KeyButtonsRes, opts KeyButtonsOpts) (comp KeyButtonsComp) {
	keyCount := len(opts.ws)
	comp.spritesList = make([][2]draws.Sprite, keyCount)
	for k := range comp.spritesList {
		for i, img := range res.imgs {
			s := draws.NewSprite(img)
			s.SetSize(opts.ws[k], opts.h)
			s.Locate(opts.xs[k], opts.y, draws.CenterTop)
			comp.spritesList[k][i] = s
		}
	}
	comp.startTimes = make([]time.Time, keyCount)
	comp.minDuration = 30
	return
}

func (comp *KeyButtonsComp) Update(keyDowns []bool) {
	comp.keyDowns = keyDowns
	for k, down := range keyDowns {
		if down {
			comp.startTimes[k] = time.Now()
		}
	}
}

// Draw key-down buttons for a while even if the press is brief.
func (comp KeyButtonsComp) Draw(dst draws.Image) {
	const (
		up   = 0
		down = 1
	)
	elapsed := time.Since(comp.startTimes[0]).Milliseconds()
	for k, keyDown := range comp.keyDowns {
		if keyDown || int32(elapsed) <= comp.minDuration {
			comp.spritesList[k][down].Draw(dst)
		} else {
			comp.spritesList[k][up].Draw(dst)
		}
	}
}
