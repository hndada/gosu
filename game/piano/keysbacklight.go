package piano

import (
	"image/color"
	"io/fs"
	"time"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/times"
)

type KeysBacklightResources struct {
	img draws.Image
}

func (br *KeysBacklightResources) Load(fsys fs.FS) {
	fname := "piano/key/backlight.png"
	br.img = draws.NewImageFromFile(fsys, fname)
}

type KeysBacklightOptions struct {
	keyCount int
	kw       []float64
	kx       []float64
	y        float64 // center bottom
	order    []KeyKind
	Colors   [4]color.NRGBA
}

func NewKeysBacklightOptions(keys KeysOptions) KeysBacklightOptions {
	return KeysBacklightOptions{
		keyCount: keys.keyCount,
		kw:       keys.kw,
		kx:       keys.kx,
		y:        keys.y,
		order:    keys.Order(),
		Colors: [4]color.NRGBA{
			{224, 224, 224, 64}, // One: white
			{255, 170, 204, 64}, // Two: pink
			{224, 224, 0, 64},   // Mid: yellow
			{255, 0, 0, 64},     // Tip: red
		},
	}
}

type KeysBacklightComponent struct {
	keysPressed []bool
	sprites     []draws.Sprite
	startTimes  []time.Time
	minDuration time.Duration
}

func NewKeysBacklightComponent(res KeysBacklightResources, opts KeysBacklightOptions) (cmp KeysBacklightComponent) {
	cmp.sprites = make([]draws.Sprite, opts.keyCount)
	for k := range cmp.sprites {
		s := draws.NewSprite(res.img)
		s.MultiplyScale(opts.kw[k] / s.W())
		s.Locate(opts.kx[k], opts.y, draws.CenterBottom)
		s.ColorScale.ScaleWithColor(opts.Colors[opts.order[k]])
		cmp.sprites[k] = s
	}
	cmp.startTimes = make([]time.Time, opts.keyCount)
	cmp.minDuration = 30 * time.Millisecond
	return
}

func (cmp *KeysBacklightComponent) Update(kp []bool) {
	cmp.keysPressed = kp
	for k, down := range kp {
		if down {
			cmp.startTimes[k] = times.Now()
		}
	}
}

// Draw backlights for a while even if the press is brief.
func (cmp KeysBacklightComponent) Draw(dst draws.Image) {
	elapsed := times.Since(cmp.startTimes[0])
	for k, p := range cmp.keysPressed {
		if p || elapsed <= cmp.minDuration {
			cmp.sprites[k].Draw(dst)
		} else {
			cmp.sprites[k].Draw(dst)
		}
	}
}
