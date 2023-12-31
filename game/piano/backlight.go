package piano

import (
	"image/color"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/times"
)

type BacklightsResources struct {
	img draws.Image
}

func (br *BacklightsResources) Load(fsys fs.FS) {
	fname := "piano/key/backlight.png"
	br.img = draws.NewImageFromFile(fsys, fname)
}

type BacklightsOptions struct {
	keyCount int
	order    []KeyKind
	keysW    []float64
	keysX    []float64
	y        float64 // center bottom
	Colors   [4]color.NRGBA
}

func NewBacklightsOptions(keys KeysOptions) BacklightsOptions {
	return BacklightsOptions{
		keyCount: keys.keyCount,
		order:    keys.Order(),
		keysW:    keys.w,
		keysX:    keys.x,
		y:        keys.y,
		Colors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
	}
}

type BacklightsComponent struct {
	sprites     []draws.Sprite
	keysPressed []bool
	startTimes  []time.Time
	minDuration time.Duration
}

func NewBacklightsComponent(res BacklightsResources, opts BacklightsOptions) (cmp BacklightsComponent) {
	cmp.sprites = make([]draws.Sprite, opts.keyCount)
	for k := range cmp.sprites {
		s := draws.NewSprite(res.img)
		s.MultiplyScale(opts.keysW[k] / s.W())
		s.Locate(opts.keysX[k], opts.y, draws.CenterBottom)
		s.ColorScale.ScaleWithColor(opts.Colors[opts.order[k]])
		cmp.sprites[k] = s
	}
	cmp.keysPressed = make([]bool, opts.keyCount)
	cmp.startTimes = make([]time.Time, opts.keyCount)
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
