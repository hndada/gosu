package piano

import (
	"time"

	draws "github.com/hndada/gosu/draws6"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/times"
)

type BacklightsComponent struct {
	sprites     []draws.Sprite
	keysPressed []bool
	startTimes  []time.Time
	minDuration time.Duration
}

func NewBacklightsComponent(res *Resources, opts *Options, keyCount int) (cmp BacklightsComponent) {
	cmp.sprites = make([]draws.Sprite, keyCount)
	ws := opts.keyWidthsMap[keyCount]
	xs := opts.keyPositionXsMap[keyCount]
	orders := opts.KeyOrders[keyCount]
	for k := range cmp.sprites {
		s := draws.NewSprite(res.BacklightsImage)
		s.Scale(ws[k] / s.W())
		s.Locate(xs[k], opts.KeyPositionY, draws.CenterBottom)
		s.ColorScale.ScaleWithColor(opts.BacklightColors[orders[k]])
		cmp.sprites[k] = s
	}
	cmp.keysPressed = make([]bool, keyCount)
	cmp.startTimes = make([]time.Time, keyCount)
	cmp.minDuration = 30 * time.Millisecond
	return
}

func (cmp *BacklightsComponent) Update(ka game.KeyboardAction) {
	kp := ka.KeysPressed()
	cmp.keysPressed = kp
	for k, p := range kp {
		if p {
			cmp.startTimes[k] = times.Now()
		}
	}
}

// Draw backlights for a while even if the press is brief.
func (cmp BacklightsComponent) Draw(dst draws.Image) {
	elapsed := times.Since(cmp.startTimes[0])
	for k, p := range cmp.keysPressed {
		if p || elapsed <= cmp.minDuration {
			cmp.sprites[k].Draw(dst)
		} else {
			cmp.sprites[k].Draw(dst)
		}
	}
}
