package piano

import (
	"time"

	draws "github.com/hndada/gosu/draws6"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/times"
)

type KeyButtonsComponent struct {
	keysSprites [][2]draws.Sprite
	keysPressed []bool
	startTimes  []time.Time
	minDuration time.Duration
}

func NewKeyButtonsComponent(res *Resources, opts *Options, keyCount int) (cmp KeyButtonsComponent) {
	cmp.keysSprites = make([][2]draws.Sprite, keyCount)
	ws := opts.keyWidthsMap[keyCount]
	for k := range cmp.keysSprites {
		for i, img := range res.KeyButtonsImages {
			s := draws.NewSprite(img)
			s.SetSize(opts.StageWidths[keyCount], opts.keyButtonHeight)
			s.Locate(ws[k], opts.KeyPositionY, draws.CenterTop)
			cmp.keysSprites[k][i] = s
		}
	}
	cmp.keysPressed = make([]bool, keyCount)
	cmp.startTimes = make([]time.Time, keyCount)
	cmp.minDuration = 30 * time.Millisecond
	return
}

func (cmp *KeyButtonsComponent) Update(ka game.KeyboardAction) {
	for k, p := range ka.KeysPressed() {
		if p {
			cmp.startTimes[k] = times.Now()
		}
	}
}

// Draw key-down buttons for a while even if the press is brief.
func (cmp KeyButtonsComponent) Draw(dst draws.Image) {
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
