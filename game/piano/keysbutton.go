package piano

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/times"
)

type KeysButtonRes struct {
	imgs [2]draws.Image
}

func (kr *KeysButtonRes) Load(fsys fs.FS) {
	for i, name := range []string{"up", "down"} {
		fname := fmt.Sprintf("piano/key/%s.png", name)
		kr.imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
}

type KeysButtonOpts struct {
	keyCount int
	kw       []float64
	h        float64
	kx       []float64
	y        float64 // center top
}

func NewKeysButtonOpts(keys KeysOpts) KeysButtonOpts {
	return KeysButtonOpts{
		keyCount: keys.keyCount,
		kw:       keys.kw,
		h:        game.ScreenH - keys.y,
		kx:       keys.kx,
		y:        keys.y,
	}
}

// Put suffix 'List' when suffix 's' is not available.
type KeysButtonComp struct {
	pressedList []bool
	spritesList [][2]draws.Sprite
	startTimes  []time.Time
	minDuration time.Duration
}

func NewKeysButtonComp(res KeysButtonRes, opts KeysButtonOpts) (comp KeysButtonComp) {
	comp.spritesList = make([][2]draws.Sprite, opts.keyCount)
	for k := range comp.spritesList {
		for i, img := range res.imgs {
			s := draws.NewSprite(img)
			s.SetSize(opts.kw[k], opts.h)
			s.Locate(opts.kx[k], opts.y, draws.CenterTop)
			comp.spritesList[k][i] = s
		}
	}
	comp.startTimes = make([]time.Time, opts.keyCount)
	comp.minDuration = 30 * time.Millisecond
	return
}

func (comp *KeysButtonComp) Update(pressedList []bool) {
	comp.pressedList = pressedList
	for k, pressed := range pressedList {
		if pressed {
			comp.startTimes[k] = times.Now()
		}
	}
}

// Draw key-down buttons for a while even if the press is brief.
func (comp KeysButtonComp) Draw(dst draws.Image) {
	const (
		up   = 0
		down = 1
	)
	elapsed := times.Since(comp.startTimes[0])
	for k, pressed := range comp.pressedList {
		if pressed || elapsed <= comp.minDuration {
			comp.spritesList[k][down].Draw(dst)
		} else {
			comp.spritesList[k][up].Draw(dst)
		}
	}
}
