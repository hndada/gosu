package piano

import (
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/times"
)

type Backlights struct {
	sprites     []draws.Sprite
	keysPressed []bool
	startTimes  []time.Time
	minDuration time.Duration
}

func NewBacklights(res *Resources, opts *Options, keyCount int) Backlights {
	backlights := Backlights{}
	backlights.sprites = make([]draws.Sprite, keyCount)
	ws := opts.keyWidthsMap[keyCount]
	xs := opts.keyPositionXsMap[keyCount]
	orders := opts.KeyOrders[keyCount]
	for k := range backlights.sprites {
		s := draws.NewSprite(res.BacklightsImage)
		s.Scale(ws[k] / s.W())
		// I thank to the past time of myself,
		// who had found the following parameters.
		s.Locate(xs[k], opts.KeyPositionY, draws.CenterBottom)
		s.ColorScale.ScaleWithColor(opts.BacklightColors[orders[k]])
		backlights.sprites[k] = s
	}
	backlights.keysPressed = make([]bool, keyCount)
	backlights.startTimes = make([]time.Time, keyCount)
	for k := range backlights.startTimes {
		backlights.startTimes[k] = times.Now()
	}
	backlights.minDuration = 30 * time.Millisecond
	return backlights
}

func (backlights *Backlights) Update(ka game.KeyboardAction) {
	kp := ka.KeysPressed()
	for k, p := range kp {
		lp := backlights.keysPressed[k]
		if (!lp && p) || (lp && !p) {
			backlights.startTimes[k] = times.Now()
		}
	}
	backlights.keysPressed = kp
}

// Draw backlights for a while even if the press is brief.
func (backlights Backlights) Draw(dst draws.Image) {
	for k, p := range backlights.keysPressed {
		elapsed := times.Since(backlights.startTimes[k])
		if p || elapsed <= backlights.minDuration {
			backlights.sprites[k].Draw(dst)
		}
	}
}
