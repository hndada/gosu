package piano

import (
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/times"
)

const (
	keyButtonUp   = 0
	keyButtonDown = 1
)

type KeyButtons struct {
	keysSprites [][2]draws.Sprite
	keysPressed []bool
	startTimes  []time.Time
	minDuration time.Duration
}

func NewKeyButtons(res *Resources, opts *Options, keyCount int) KeyButtons {
	keyButtons := KeyButtons{}
	keyButtons.keysSprites = make([][2]draws.Sprite, keyCount)
	ws := opts.keyWidthsMap[keyCount]
	xs := opts.keyPositionXsMap[keyCount]
	for k := range keyButtons.keysSprites {
		for i, img := range res.KeyButtonsImages {
			s := draws.NewSprite(img)
			s.SetSize(ws[k], opts.keyButtonHeight)
			s.Locate(xs[k], opts.KeyPositionY, draws.CenterTop)
			keyButtons.keysSprites[k][i] = s
		}
	}
	keyButtons.keysPressed = make([]bool, keyCount)
	keyButtons.startTimes = make([]time.Time, keyCount)
	for k := range keyButtons.startTimes {
		keyButtons.startTimes[k] = times.Now()
	}
	keyButtons.minDuration = 30 * time.Millisecond
	return keyButtons
}

func (keyButtons *KeyButtons) Update(ka game.KeyboardAction) {
	for k, p := range ka.KeysPressed() {
		if p {
			keyButtons.startTimes[k] = times.Now()
		}
	}
}

// Draw key-down buttons for a while even if the press is brief.
func (keyButtons KeyButtons) Draw(dst draws.Image) {
	for k, p := range keyButtons.keysPressed {
		elapsed := times.Since(keyButtons.startTimes[k])
		if p || elapsed <= keyButtons.minDuration {
			keyButtons.keysSprites[k][keyButtonDown].Draw(dst)
		} else {
			keyButtons.keysSprites[k][keyButtonUp].Draw(dst)
		}
	}
}
