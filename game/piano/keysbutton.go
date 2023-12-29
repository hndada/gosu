package piano

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/times"
)

type KeysButtonResources struct {
	imgs [2]draws.Image
}

func (kr *KeysButtonResources) Load(fsys fs.FS) {
	for i, name := range []string{"up", "down"} {
		fname := fmt.Sprintf("piano/key/%s.png", name)
		kr.imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
}

type KeysButtonOptions struct {
	keyCount int
	kw       []float64
	h        float64
	kx       []float64
	y        float64 // center top
}

func NewKeysButtonOptions(keys KeysOptions) KeysButtonOptions {
	return KeysButtonOptions{
		keyCount: keys.keyCount,
		kw:       keys.kw,
		h:        game.ScreenH - keys.y,
		kx:       keys.kx,
		y:        keys.y,
	}
}

type KeysButtonComponent struct {
	keysPressed []bool
	keysSprites [][2]draws.Sprite
	startTimes  []time.Time
	minDuration time.Duration
}

func NewKeysButtonComponent(res KeysButtonResources, opts KeysButtonOptions) (cmp KeysButtonComponent) {
	cmp.keysSprites = make([][2]draws.Sprite, opts.keyCount)
	for k := range cmp.keysSprites {
		for i, img := range res.imgs {
			s := draws.NewSprite(img)
			s.SetSize(opts.kw[k], opts.h)
			s.Locate(opts.kx[k], opts.y, draws.CenterTop)
			cmp.keysSprites[k][i] = s
		}
	}
	cmp.startTimes = make([]time.Time, opts.keyCount)
	cmp.minDuration = 30 * time.Millisecond
	return
}

func (cmp *KeysButtonComponent) Update(kp []bool) {
	cmp.keysPressed = kp
	for k, p := range kp {
		if p {
			cmp.startTimes[k] = times.Now()
		}
	}
}

// Draw key-down buttons for a while even if the press is brief.
func (cmp KeysButtonComponent) Draw(dst draws.Image) {
	const (
		up   = 0
		down = 1
	)
	elapsed := times.Since(cmp.startTimes[0])
	for k, p := range cmp.keysPressed {
		if p || elapsed <= cmp.minDuration {
			cmp.keysSprites[k][down].Draw(dst)
		} else {
			cmp.keysSprites[k][up].Draw(dst)
		}
	}
}
